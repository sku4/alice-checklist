package googlekeep

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	ClientPlatform        = "WEB"
	ClientLocale          = "ru"
	ClientVersionBuild    = "0"
	ClientVersionMajor    = "3"
	ClientVersionMinor    = "3"
	ClientVersionRevision = "167"
	NodeBaseVersion       = "0"
	NodeKind              = "notes#node"
	NodeTimestampsKind    = "notes#timestamps"
)

var (
	HeaderCapabilities = []string{
		"EC", "TR", "SH", "LB", "RB", "DR", "AN", "PI", "EX", "IN", "SNB", "CO", "MI", "NC", "IN",
	}
)

type Request struct {
	TargetVersion   string    `json:"targetVersion"`
	ClientTimestamp time.Time `json:"clientTimestamp"`
	Nodes           []Node    `json:"nodes"`
	RequestHeader   `json:"requestHeader"`
}

func NewRequest() (pd *Request) {
	pd = &Request{}
	pd.RequestHeader = *NewRequestHeader()
	pd.ClientTimestamp = time.Now()
	return
}

type Node struct {
	Id            string `json:"id"`
	Kind          string `json:"kind"`
	ParentId      string `json:"parentId"`
	Timestamps    `json:"timestamps"`
	Type          string `json:"type"`
	TrashState    int    `json:"trashState"`
	ServerId      string `json:"serverId,omitempty"`
	DeletionState int    `json:"deletionState"`
	SortValue     string `json:"sortValue"`
	BaseVersion   string `json:"baseVersion"`
	Title         string `json:"title,omitempty"`
	IsArchived    bool   `json:"isArchived"`
	IsPinned      bool   `json:"isPinned"`
	Background    struct {
		Name   string `json:"name,omitempty"`
		Origin string `json:"origin,omitempty"`
	} `json:"background,omitempty"`
	NodeSettings struct {
		GraveyardState string `json:"graveyardState,omitempty"`
	} `json:"nodeSettings,omitempty"`
	ParentServerId        string `json:"parentServerId,omitempty"`
	Text                  string `json:"text,omitempty"`
	Checked               bool   `json:"checked,omitempty"`
	SuperListItemId       string `json:"superListItemId,omitempty"`
	SuperListItemServerId string `json:"superListItemServerId,omitempty"`
}

func NewNode() (n *Node) {
	n = &Node{}
	n.BaseVersion = NodeBaseVersion
	n.Kind = NodeKind
	n.Timestamps = Timestamps{
		Kind:                    NodeTimestampsKind,
		Created:                 time.Unix(0, 0),
		Deleted:                 time.Unix(0, 0),
		Updated:                 time.Unix(0, 0),
		Trashed:                 time.Unix(0, 0),
		UserEdited:              time.Unix(0, 0),
		RecentSharedChangesSeen: time.Unix(0, 0),
	}
	n.Type = "LIST_ITEM"
	return
}

func (n *Node) Add(parentId, text string) {
	n.clean()
	n.Timestamps.Created = time.Now()
	n.ParentId = parentId
	n.Text = text
	n.SortValue = strconv.FormatInt(time.Now().UnixMicro()+genSortId(), 10)
	n.genId()
}

func (n *Node) Update(checked bool) {
	n.clean()
	n.Timestamps.Updated = time.Now()
	n.Checked = checked
	n.DeletionState = 0
	n.SortValue = strconv.FormatInt(time.Now().UnixMicro()+genSortId(), 10)
}

func (n *Node) Delete() {
	n.clean()
	n.Timestamps.Deleted = time.Now()
	n.DeletionState = 1
}

func (n *Node) clean() {
	n.Timestamps.Kind = NodeTimestampsKind
	if n.Timestamps.Created.Unix() <= 0 {
		n.Timestamps.Created = time.Unix(0, 0)
	}
	if n.Timestamps.Deleted.Unix() <= 0 {
		n.Timestamps.Deleted = time.Unix(0, 0)
	}
	if n.Timestamps.Updated.Unix() <= 0 {
		n.Timestamps.Updated = time.Unix(0, 0)
	}
	if n.Timestamps.Trashed.Unix() <= 0 {
		n.Timestamps.Trashed = time.Unix(0, 0)
	}
	if n.Timestamps.UserEdited.Unix() <= 0 {
		n.Timestamps.UserEdited = time.Unix(0, 0)
	}
	if n.Timestamps.RecentSharedChangesSeen.Unix() <= 0 {
		n.Timestamps.RecentSharedChangesSeen = time.Unix(0, 0)
	}
}

var genRandId = func() func() string {
	c := 0
	return func() string {
		i := []rune(strconv.FormatInt(time.Now().UnixNano(), 10))
		s1 := rand.NewSource(time.Now().UnixNano() + int64(c))
		r1 := rand.New(s1)
		r2 := r1.Intn(88888888) + 11111111
		c++
		return string(i[:13]) + "." + strconv.FormatInt(int64(r2), 10)
	}
}()

func (n *Node) genId() {
	n.Id = genRandId()
}

var genSortId = func() func() int64 {
	var c int64 = -1
	return func() int64 {
		c++
		return c
	}
}()

type Timestamps struct {
	Kind                    string    `json:"kind"`
	Created                 time.Time `json:"created"`
	Deleted                 time.Time `json:"deleted"`
	Trashed                 time.Time `json:"trashed"`
	Updated                 time.Time `json:"updated"`
	UserEdited              time.Time `json:"userEdited"`
	RecentSharedChangesSeen time.Time `json:"recentSharedChangesSeen,omitempty"`
}

type RequestHeader struct {
	RequestId       string `json:"requestId,omitempty"`
	ClientVersion   `json:"clientVersion"`
	ClientPlatform  string       `json:"clientPlatform"`
	Capabilities    []Capability `json:"capabilities"`
	ClientSessionId string       `json:"clientSessionId"`
	ClientLocale    string       `json:"clientLocale"`
}

func NewRequestHeader() (rh *RequestHeader) {
	rh = &RequestHeader{}
	rh.ClientPlatform = ClientPlatform
	rh.ClientLocale = ClientLocale
	rh.ClientVersion.Build = ClientVersionBuild
	rh.ClientVersion.Major = ClientVersionMajor
	rh.ClientVersion.Minor = ClientVersionMinor
	rh.ClientVersion.Revision = ClientVersionRevision
	for _, capability := range HeaderCapabilities {
		rh.Capabilities = append(rh.Capabilities, Capability{
			Type: capability,
		})
	}
	return
}

type Capability struct {
	Type string `json:"type"`
}

type ClientVersion struct {
	Major    string `json:"major"`
	Minor    string `json:"minor"`
	Build    string `json:"build"`
	Revision string `json:"revision"`
}
