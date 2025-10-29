// LL2Launch represents the launch data structure from Launch Library 2 API
// https://ll.thespacedevs.com/2.3.0/json
package models

type LL2Response struct {
	Count    int                  `json:"count"`
	Next     string               `json:"next"`
	Previous string               `json:"previous"`
	Results  []*LL2LaunchDetailed `json:"results"`
}

type LL2LaunchBasic struct {
	ID               string          `json:"id"`
	URL              string          `json:"url"`
	Name             string          `json:"name"`
	ResponseMode     string          `json:"response_mode"`
	Slug             string          `json:"slug"`
	LaunchDesignator string          `json:"launch_designator"`
	Status           LL2Status       `json:"status"`
	LastUpdated      string          `json:"last_updated"`
	Net              string          `json:"net"`
	NetPrecision     LL2NetPrecision `json:"net_precision"`
	WindowEnd        string          `json:"window_end"`
	WindowStart      string          `json:"window_start"`
	Image            LL2Image        `json:"image"`
	Infographic      string          `json:"infographic"`
}

type LL2Image struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	ImageURL     string            `json:"image_url"`
	ThumbnailURL string            `json:"thumbnail_url"`
	Credit       string            `json:"credit,omitempty"`
	License      LL2ImageLicense   `json:"license"`
	SingleUse    bool              `json:"single_use"`
	Variants     []LL2ImageVariant `json:"variants"`
}

type LL2ImageLicense struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Link     string `json:"link"`
}

type LL2ImageVariant struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type LL2Status struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type LL2NetPrecision struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type LL2LaunchNormal struct {
	LL2LaunchBasic
	Probability                    int                `json:"probability"`
	WeatherConcerns                string             `json:"weather_concerns"`
	FailReason                     string             `json:"failreason"`
	HashTag                        string             `json:"hashtag"`
	LaunchServiceProvider          LL2AgencyMini      `json:"launch_service_provider"`
	Rocket                         LL2RocketNormal    `json:"rocket"`
	Mission                        LL2Mission         `json:"mission"`
	Pad                            LL2Pad             `json:"pad"`
	WebcastLive                    bool               `json:"webcast_live"`
	Program                        []LL2ProgramNormal `json:"program"`
	OrbitalLaunchAttemptCount      int                `json:"orbital_launch_attempt_count"`
	LocationLaunchAttemptCount     int                `json:"location_launch_attempt_count"`
	PadLaunchAttemptCount          int                `json:"pad_launch_attempt_count"`
	AgencyLaunchAttemptCount       int                `json:"agency_launch_attempt_count"`
	OrbitalLaunchAttemptCountYear  int                `json:"orbital_launch_attempt_count_year"`
	LocationLaunchAttemptCountYear int                `json:"location_launch_attempt_count_year"`
	PadLaunchAttemptCountYear      int                `json:"pad_launch_attempt_count_year"`
	AgencyLaunchAttemptCountYear   int                `json:"agency_launch_attempt_count_year"`
}

type LL2AgencyMini struct {
	ResponseMode string        `json:"response_mode"`
	ID           int           `json:"id"`
	URL          string        `json:"url"`
	Name         string        `json:"name"`
	Abbrev       string        `json:"abbrev"`
	Type         LL2AgencyType `json:"type"`
}

type LL2AgencyType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LL2RocketNormal struct {
	ID            int                   `json:"id"`
	Configuration LL2LauncherConfigList `json:"configuration"`
}

type LL2LauncherConfigList struct {
	ResponseMode string                        `json:"response_mode"`
	ID           int                           `json:"id"`
	URL          string                        `json:"url"`
	Name         string                        `json:"name"`
	Families     []LL2LauncherConfigFamilyMini `json:"families"`
	FullName     string                        `json:"full_name"`
	Variant      string                        `json:"variant"`
}

type LL2LauncherConfigFamilyMini struct {
	ResponseMode string `json:"response_mode"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
}

type LL2Mission struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
	Image       LL2Image            `json:"image"`
	Orbit       LL2Orbit            `json:"orbit"`
	Agencies    []LL2AgencyDetailed `json:"agencies"`
	InfoURLs    []LL2InfoURL        `json:"info_urls"`
	VidURLs     []LL2VidURL         `json:"vid_urls"`
}

type LL2Orbit struct {
	ID            int                  `json:"id"`
	Name          string               `json:"name"`
	Abbrev        string               `json:"abbrev"`
	CelestialBody LL2CelestialBodyMini `json:"celestial_body"`
}

type LL2CelestialBodyMini struct {
	ResponseMode string `json:"response_mode"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
}

type LL2AgencyDetailed struct {
	ResponseMode                  string               `json:"response_mode"`
	ID                            int                  `json:"id"`
	URL                           string               `json:"url"`
	Name                          string               `json:"name"`
	Abbrev                        string               `json:"abbrev"`
	Type                          LL2AgencyType        `json:"type"`
	Featured                      bool                 `json:"featured"`
	Country                       []LL2Country         `json:"country"`
	Description                   string               `json:"description"`
	Administrator                 string               `json:"administrator"`
	FoundingYear                  int                  `json:"founding_year"`
	Launchers                     string               `json:"launchers"`
	Spacecraft                    string               `json:"spacecraft"`
	Parent                        string               `json:"parent"`
	Image                         LL2Image             `json:"image"`
	Logo                          LL2Image             `json:"logo"`
	SocialLogo                    LL2Image             `json:"social_logo"`
	TotalLaunchCount              int                  `json:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int                  `json:"consecutive_successful_launches"`
	SuccessfulLaunches            int                  `json:"successful_launches"`
	FailedLaunches                int                  `json:"failed_launches"`
	PendingLaunches               int                  `json:"pending_launches"`
	ConsecutiveSuccessfulLandings int                  `json:"consecutive_successful_landings"`
	SuccessfulLandings            int                  `json:"successful_landings"`
	FailedLandings                int                  `json:"failed_landings"`
	AttemptedLandings             int                  `json:"attempted_landings"`
	SuccessfulLandingsSpacecraft  int                  `json:"successful_landings_spacecraft"`
	FailedLandingsSpacecraft      int                  `json:"failed_landings_spacecraft"`
	AttemptedLandingsSpacecraft   int                  `json:"attempted_landings_spacecraft"`
	SuccessfulLandingsPayload     int                  `json:"successful_landings_payload"`
	FailedLandingsPayload         int                  `json:"failed_landings_payload"`
	AttemptedLandingsPayload      int                  `json:"attempted_landings_payload"`
	InfoURL                       string               `json:"info_url"`
	WikiURL                       string               `json:"wiki_url"`
	SocialMediaLinks              []LL2SocialMediaLink `json:"social_media_links"`
}

type LL2Country struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	Alpha2Code              string `json:"alpha2_code"`
	Alpha3Code              string `json:"alpha3_code"`
	NationalityName         string `json:"nationality_name"`
	NationalityNameComposed string `json:"nationality_name_composed"`
}

type LL2SocialMediaLink struct {
	ID          int            `json:"id"`
	SocialMedia LL2SocialMedia `json:"social_media"`
	URL         string         `json:"url"`
}

type LL2SocialMedia struct {
	Id   int      `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Logo LL2Image `json:"logo"`
}

type LL2InfoURL struct {
	Priority     int            `json:"priority"`
	Source       string         `json:"source"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	FeatureImage string         `json:"feature_image"`
	URL          string         `json:"url"`
	Type         LL2InfoURLType `json:"type"`
	Language     LL2Language    `json:"language"`
}

type LL2InfoURLType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type LL2Language struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type LL2VidURL struct {
	Priority     int           `json:"priority"`
	Source       string        `json:"source"`
	Publisher    string        `json:"publisher"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	FeatureImage string        `json:"feature_image"`
	URL          string        `json:"url"`
	Type         LL2VidURLType `json:"type"`
	Language     LL2Language   `json:"language"`
	StartTime    string        `json:"start_time"`
	EndTime      string        `json:"end_time"`
	Live         bool          `json:"live"`
}

type LL2VidURLType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type LL2Pad struct {
	Id                        int               `json:"id"`
	URL                       string            `json:"url"`
	Active                    bool              `json:"active"`
	Agencies                  []LL2AgencyNormal `json:"agencies"`
	Name                      string            `json:"name"`
	Image                     LL2Image          `json:"image"`
	Description               string            `json:"description"`
	InfoURL                   string            `json:"info_url"`
	WikiURL                   string            `json:"wiki_url"`
	MapURL                    string            `json:"map_url"`
	Latitude                  float64           `json:"latitude"`
	Longitude                 float64           `json:"longitude"`
	Country                   LL2Country        `json:"country"`
	MapImage                  string            `json:"map_image"`
	TotalLaunchCount          int               `json:"total_launch_count"`
	OrbitalLaunchAttemptCount int               `json:"orbital_launch_attempt_count"`
	FastestTurnaround         string            `json:"fastest_turnaround"`
	Location                  LL2Location       `json:"location"`
}

type LL2Location struct {
	ResponseMode      string                   `json:"response_mode"`
	ID                int                      `json:"id"`
	Name              string                   `json:"name"`
	URL               string                   `json:"url"`
	CelestialBody     LL2CelestialBodyDetailed `json:"celestial_body"`
	Active            bool                     `json:"active"`
	Country           LL2Country               `json:"country"`
	Description       string                   `json:"description"`
	Image             LL2Image                 `json:"image"`
	MapImage          string                   `json:"map_image"`
	Latitude          float64                  `json:"latitude"`
	Longitude         float64                  `json:"longitude"`
	TimezoneName      string                   `json:"timezone_name"`
	TotalLaunchCount  int                      `json:"total_launch_count"`
	TotalLandingCount int                      `json:"total_landing_count"`
}

type LL2CelestialBodyDetailed struct {
	LL2CelestialBodyMini
	Type                   LL2CelestialBodyType `json:"type"`
	Diameter               float64              `json:"diameter"`
	Mass                   float64              `json:"mass"`
	Gravity                float64              `json:"gravity"`
	LengthOfDay            string               `json:"length_of_day"`
	Atmosphere             bool                 `json:"atmosphere"`
	Image                  LL2Image             `json:"image"`
	Description            string               `json:"description"`
	WikiURL                string               `json:"wiki_url"`
	TotalAttemptedLandes   int                  `json:"total_attempted_landes"`
	SuccessfulLaunches     int                  `json:"successful_launches"`
	FailedLaunches         int                  `json:"failed_launches"`
	TotalAttemptedLandings int                  `json:"total_attempted_landings"`
	SuccessfulLandings     int                  `json:"successful_landings"`
	FailedLandings         int                  `json:"failed_landings"`
}

type LL2CelestialBodyType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LL2AgencyNormal struct {
	LL2AgencyMini
	Featured      bool         `json:"featured"`
	Country       []LL2Country `json:"country"`
	Description   string       `json:"description"`
	Administrator string       `json:"administrator"`
	FoundingYear  int          `json:"founding_year"`
	Launchers     string       `json:"launchers"`
	Spacecraft    string       `json:"spacecraft"`
	Parent        string       `json:"parent"`
	Image         LL2Image     `json:"image"`
	Logo          LL2Image     `json:"logo"`
	SocialLogo    LL2Image     `json:"social_logo"`
}

type LL2ProgramNormal struct {
	ResponseMode   string            `json:"response_mode"`
	ID             int               `json:"id"`
	URL            string            `json:"url"`
	Name           string            `json:"name"`
	Image          LL2Image          `json:"image"`
	InfoUrl        string            `json:"info_url"`
	WikiUrl        string            `json:"wiki_url"`
	Description    string            `json:"description"`
	Agencies       []LL2AgencyMini   `json:"agencies"`
	StartDate      string            `json:"start_date"`
	EndDate        string            `json:"end_date"`
	MissionPatches []LL2MissionPatch `json:"mission_patches"`
	Type           LL2ProgramType    `json:"type"`
}
type LL2MissionPatch struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Priority     int           `json:"priority"`
	ImageUrl     string        `json:"image_url"`
	Agency       LL2AgencyMini `json:"agency"`
	ResponseMode string        `json:"response_mode"`
}

type LL2ProgramType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LL2LaunchDetailed struct {
	LL2LaunchNormal
	FlightclubUrl  string             `json:"flightclub_url"`
	Updates        []LL2Update        `json:"updates"`
	InfoURLs       []LL2InfoURL       `json:"info_urls"`
	VidURLs        []LL2VidURL        `json:"vid_urls"`
	Timeline       []LL2TimelineEvent `json:"timeline"`
	PadTurnaround  string             `json:"pad_turnaround"`
	MissionPatches []LL2MissionPatch  `json:"mission_patches"`
}

type LL2Update struct {
	ID           int    `json:"id"`
	ProfileImage string `json:"profile_image"`
	Comment      string `json:"comment"`
	InfoUrl      string `json:"info_url"`
	CreatedBy    string `json:"created_by"`
	CreatedOn    string `json:"created_on"`
}

type LL2TimelineEvent struct {
	Type         LL2TimelineEventType `json:"type"`
	RelativeTime int                  `json:"relative_time"`
}

type LL2TimelineEventType struct {
	ID          int    `json:"id"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}
