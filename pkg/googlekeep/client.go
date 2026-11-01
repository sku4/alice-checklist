package googlekeep

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/har"
	"github.com/sirupsen/logrus"
	"github.com/sku4/alice-checklist/configs"
	"github.com/sku4/alice-checklist/lang"
	"github.com/sku4/alice-checklist/models/googlekeep"
	"github.com/sku4/alice-checklist/pkg/boltdb"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	FileHarName            = "./configs/googlekeep/keep.google.com.har"
	HeaderXOriginKey       = "x-origin"
	HeaderAuthorizationKey = "authorization"
	CookieContainPrefix    = "SID"
	UrlContainPrefix1      = "notes-pa"
	UrlContainPrefix2      = "changes"
	UrlContainPrefix3      = "google.com"
	// CookieSapisidKey holds the cookie whose value is used to compute SAPISIDHASH.
	CookieSapisidKey = "SAPISID"
	// CookieSapisidFallbackKey is used when SAPISID is missing (e.g. Secure-only contexts).
	CookieSapisidFallbackKey = "__Secure-3PAPISID"
	// OriginKeep is the default origin used for SAPISIDHASH computation.
	OriginKeep        = "https://keep.google.com"
	SapisidHashPrefix = "SAPISIDHASH "
)

type Client struct {
	cache  boltdb.CacheRepository
	config *configs.Config
	loc    *lang.Localize
}

func NewClient(loc *lang.Localize, cfg *configs.Config, db boltdb.Storage) *Client {
	return &Client{
		cache:  boltdb.NewNodeCacheRepository(db),
		config: cfg,
		loc:    loc,
	}
}

func (r *Client) Patch(add, delete []string) (err error) {
	// if delete items from mobile, then need force reload notes
	//_, err = r.reloadNotes()

	postData := googlekeep.NewRequest()
	noteRootId := r.config.GoogleKeep.NoteRootId

	needReload := false
	for _, s := range delete {
		_, err := r.cache.Get(s)
		if err == boltdb.ErrorClearTextEmpty {
			continue
		} else if err != nil {
			needReload = true
			break
		}
	}

	if needReload {
		_, err = r.reloadNotes()
		if err != nil {
			return
		}
	}

	for _, s := range delete {
		entryNode, err := r.cache.Get(s)
		if err == boltdb.ErrorClearTextEmpty {
			continue
		} else if err != nil && err != boltdb.ErrorNodeNotFound {
			return err
		}
		if err == nil {
			entryNode.Update(true)
			postData.Nodes = append(postData.Nodes, *entryNode)
		}
	}

	for _, s := range add {
		entryNode, err := r.cache.Get(s)
		if err == boltdb.ErrorClearTextEmpty {
			continue
		} else if err != nil && err != boltdb.ErrorNodeNotFound {
			return err
		}
		if err == boltdb.ErrorNodeNotFound {
			entryNode = googlekeep.NewNode()
			entryNode.Add(noteRootId, s)
		} else {
			entryNode.Update(false)
		}
		postData.Nodes = append(postData.Nodes, *entryNode)
	}

	_, err = r.send(*postData)
	return
}

func (r *Client) List() (nodes []googlekeep.Node, err error) {
	nodes, err = r.reloadNotes()
	if err != nil {
		return
	}

	return
}

func (r *Client) CacheList() (nodes []googlekeep.Node, err error) {
	nodes, err = r.cache.List()
	if err != nil {
		return
	}

	return
}

func (r *Client) Clean() (err error) {
	nodes, err := r.reloadNotes()
	if err != nil {
		return
	}

	postData := googlekeep.NewRequest()
	duplicates := make(map[string][]googlekeep.Node)
	for _, n := range nodes {
		s, err := boltdb.ClearText(n.Text)
		if err == boltdb.ErrorClearTextEmpty {
			continue
		}
		if len(duplicates[s]) == 0 {
			duplicates[s] = make([]googlekeep.Node, 0)
		}
		duplicates[s] = append(duplicates[s], n)
	}

	for _, nodes := range duplicates {
		if len(nodes) > 1 {
			checked := true
			var nodesDelete []googlekeep.Node
			for _, n := range nodes {
				if !checked || n.Checked {
					nodesDelete = append(nodesDelete, n)
				}
				if !n.Checked {
					checked = false
				}
			}
			if checked {
				for _, n := range nodes[1:] {
					n.Delete()
					postData.Nodes = append(postData.Nodes, n)
				}
			} else {
				for _, n := range nodesDelete {
					n.Delete()
					postData.Nodes = append(postData.Nodes, n)
				}
			}
		}
	}

	if len(postData.Nodes) > 0 {
		_, err = r.send(*postData)
		_, err = r.reloadNotes()
	}

	return
}

func (r *Client) checkNotes() (err error) {
	empty, err := r.cache.IsEmptyBucket()
	if err != nil {
		return
	}
	if empty {
		_, err = r.reloadNotes()
	}
	return
}

func (r *Client) reloadNotes() (nodes []googlekeep.Node, err error) {
	err = r.cache.Truncate()
	if err != nil {
		return
	}

	noteRootId := r.config.GoogleKeep.NoteRootId
	noteRootExists := false
	endOfNodes := false
	targetVersion := ""

	for !noteRootExists && !endOfNodes {
		postData := googlekeep.NewRequest()
		postData.TargetVersion = targetVersion
		entryNoteResp, err := r.send(*postData)
		if err != nil {
			return nil, err
		}

		for _, n := range entryNoteResp.Nodes {
			if n.Id == noteRootId {
				noteRootExists = true
			} else if n.ParentId == noteRootId {
				nodes = append(nodes, n)
			}
		}

		if len(entryNoteResp.Nodes) == 0 {
			endOfNodes = true
		}
		targetVersion = entryNoteResp.ToVersion
	}

	for _, n := range nodes {
		err = r.cache.Save(n)
		if err != nil {
			return
		}
	}

	return
}

func (r *Client) send(postData googlekeep.Request) (entryNoteResp googlekeep.Response, err error) {
	var h har.HAR
	dataHar, err := ioutil.ReadFile(FileHarName)
	if err != nil {
		return
	}
	if err = h.UnmarshalJSON(dataHar); err != nil {
		return
	}

	var harEntryNote *har.Entry
	for _, entry := range h.Log.Entries {
		if strings.Contains(entry.Request.URL, UrlContainPrefix1) &&
			strings.Contains(entry.Request.URL, UrlContainPrefix2) &&
			strings.Contains(entry.Request.URL, UrlContainPrefix3) {
			harEntryNote = entry
			break
		}
	}

	if harEntryNote == nil {
		return entryNoteResp, errors.New(fmt.Sprintf(
			r.loc.Translate("Not found entry note"),
		))
	}

	postDataJson, err := json.Marshal(postData)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, harEntryNote.Request.URL,
		bytes.NewBuffer(postDataJson))
	if err != nil {
		return entryNoteResp, errors.New(fmt.Sprintf(
			r.loc.Translate("Got error %s"), err.Error()))
	}

	for _, cookie := range harEntryNote.Request.Cookies {
		if strings.Contains(cookie.Name, CookieContainPrefix) {
			req.AddCookie(&http.Cookie{
				Name:   cookie.Name,
				Value:  cookie.Value,
				Domain: cookie.Domain,
				Path:   cookie.Path,
			})
		}
	}

	// SAPISID value used to recompute the SAPISIDHASH Authorization header with a
	// fresh timestamp on every request. The hash captured in the HAR carries a
	// stale timestamp that Google rejects (401 CREDENTIALS_MISSING); the SAPISID
	// cookie itself stays valid for months, so recomputing the hash keeps auth alive.
	sapisid, origin := sapisidForHash(harEntryNote)

	for _, header := range harEntryNote.Request.Headers {
		switch header.Name {
		case HeaderXOriginKey:
			req.Header.Add(header.Name, header.Value)
		case HeaderAuthorizationKey:
			if sapisid != "" && strings.HasPrefix(header.Value, SapisidHashPrefix) {
				req.Header.Add(header.Name, buildSapisidHash(sapisid, origin))
			} else {
				req.Header.Add(header.Name, header.Value)
			}
		}
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return entryNoteResp, errors.New(fmt.Sprintf(
			r.loc.Translate("Error occured. Error is: %s"), err.Error()))
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bodyString := ""
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		bodyString += scanner.Text()
	}

	if err = json.Unmarshal([]byte(bodyString), &entryNoteResp); err != nil {
		logrus.Println(bodyString)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Println(bodyString)
		return entryNoteResp, errors.New(fmt.Sprint(
			r.loc.Translate("Google Keep response status: "), resp.Status))
	}

	return
}

// sapisidForHash extracts the SAPISID cookie value (falling back to its secure
// variant) and the request origin to use for SAPISIDHASH computation.
func sapisidForHash(entry *har.Entry) (sapisid, origin string) {
	origin = OriginKeep
	for _, cookie := range entry.Request.Cookies {
		if cookie.Name == CookieSapisidKey {
			sapisid = cookie.Value
		} else if sapisid == "" && cookie.Name == CookieSapisidFallbackKey {
			sapisid = cookie.Value
		}
	}
	for _, header := range entry.Request.Headers {
		if header.Name == HeaderXOriginKey && header.Value != "" {
			origin = header.Value
		}
	}
	return
}

// buildSapisidHash computes a fresh SAPISIDHASH: "<unix_ts>_<sha1(ts SAPISID origin)>".
// Google validates the hash and rejects timestamps too far from now, so the
// timestamp is taken at call time rather than reused from the HAR.
func buildSapisidHash(sapisid, origin string) string {
	ts := time.Now().Unix()
	sum := sha1.Sum([]byte(fmt.Sprintf("%d %s %s", ts, sapisid, origin)))
	return fmt.Sprintf("%s%d_%s", SapisidHashPrefix, ts, hex.EncodeToString(sum[:]))
}
