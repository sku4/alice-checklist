package googlekeep

import (
	"time"
)

type Response struct {
	Kind      string `json:"kind"`
	ToVersion string `json:"toVersion"`
	UserInfo  struct {
		Timestamps struct {
			Kind    string    `json:"kind"`
			Created time.Time `json:"created"`
		} `json:"timestamps"`
		Settings struct {
			SingleSettings []struct {
				Type                              string   `json:"type"`
				ApplicablePlatforms               []string `json:"applicablePlatforms"`
				GlobalCheckedListItemsPolicyValue string   `json:"globalCheckedListItemsPolicyValue,omitempty"`
				GlobalNewListItemPlacementValue   string   `json:"globalNewListItemPlacementValue,omitempty"`
				LayoutStyleValue                  string   `json:"layoutStyleValue,omitempty"`
				SharingEnabledValue               bool     `json:"sharingEnabledValue,omitempty"`
				WebEmbedsEnabledValue             bool     `json:"webEmbedsEnabledValue,omitempty"`
				WebAppThemeValue                  string   `json:"webAppThemeValue,omitempty"`
			} `json:"singleSettings"`
		} `json:"settings"`
		ContextualCoachmarksAcked []string `json:"contextualCoachmarksAcked"`
	} `json:"userInfo"`
	Nodes              []Node `json:"nodes"`
	Truncated          bool   `json:"truncated"`
	UpgradeRecommended bool   `json:"upgradeRecommended"`
	ForceFullResync    bool   `json:"forceFullResync"`
	ResponseHeader     struct {
		UpdateState   string `json:"updateState"`
		RequestId     string `json:"requestId"`
		ExperimentIds []int  `json:"experimentIds"`
	} `json:"responseHeader"`
}
