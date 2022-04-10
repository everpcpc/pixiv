package pixiv

type UserImages struct {
	Medium string `json:"medium"`
}

type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Comment    string `json:"comment"`
	IsFollowed bool   `json:"is_followed"`

	ProfileImages *UserImages `json:"profile_image_urls"`
}

type UserDetail struct {
	User             *User             `json:"user"`
	Profile          *Profile          `json:"profile"`
	ProfilePublicity *ProfilePublicity `json:"profile_publicity"`
	Workspace        *Workspace        `json:"workspace"`
}

type Profile struct {
	Webpage                    interface{} `json:"webpage"`
	Gender                     string      `json:"gender"`
	Birth                      string      `json:"birth"`
	BirthDay                   string      `json:"birth_day"`
	BirthYear                  uint64      `json:"birth_year"`
	Region                     string      `json:"region"`
	AddressID                  uint64      `json:"address_id"`
	CountryCode                string      `json:"country_code"`
	Job                        string      `json:"job"`
	JobID                      uint64      `json:"job_id"`
	TotalFollowUsers           uint64      `json:"total_follow_users"`
	TotalMypixivUsers          uint64      `json:"total_mypixiv_users"`
	TotalIllusts               uint64      `json:"total_illusts"`
	TotalManga                 uint64      `json:"total_manga"`
	TotalNovels                uint64      `json:"total_novels"`
	TotalIllustBookmarksPublic uint64      `json:"total_illust_bookmarks_public"`
	TotalIllustSeries          uint64      `json:"total_illust_series"`
	TotalNovelSeries           uint64      `json:"total_novel_series"`
	BackgroundImageURL         string      `json:"background_image_url"`
	TwitterAccount             string      `json:"twitter_account"`
	TwitterURL                 string      `json:"twitter_url"`
	PawooURL                   string      `json:"pawoo_url"`
	IsPremium                  bool        `json:"is_premium"`
	IsUsingCustomProfileImage  bool        `json:"is_using_custom_profile_image"`
}

type Workspace struct {
	Pc                string `json:"pc"`
	Monitor           string `json:"monitor"`
	Tool              string `json:"tool"`
	Scanner           string `json:"scanner"`
	Tablet            string `json:"tablet"`
	Mouse             string `json:"mouse"`
	Printer           string `json:"printer"`
	Desktop           string `json:"desktop"`
	Music             string `json:"music"`
	Desk              string `json:"desk"`
	Chair             string `json:"chair"`
	Comment           string `json:"comment"`
	WorkspaceImageURL string `json:"workspace_image_url"`
}

type ProfilePublicity struct {
	Gender    string `json:"gender"`
	Region    string `json:"region"`
	BirthDay  string `json:"birth_day"`
	BirthYear string `json:"birth_year"`
	Job       string `json:"job"`
	Pawoo     bool   `json:"pawoo"`
}

type Tag struct {
	Name           string `json:"name"`
	TranslatedName string `json:"translated_name"`
}
type Images struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
	Original     string `json:"original"`
}
type MetaSinglePage struct {
	OriginalImageURL string `json:"original_image_url"`
}
type MetaPage struct {
	Images Images `json:"image_urls"`
}

type Illust struct {
	ID             uint64          `json:"id"`
	Title          string          `json:"title"`
	Type           string          `json:"type"`
	Images         *Images         `json:"image_urls"`
	Caption        string          `json:"caption"`
	Restrict       int             `json:"restrict"`
	User           *User           `json:"user"`
	Tags           []Tag           `json:"tags"`
	Tools          []string        `json:"tools"`
	CreateData     string          `json:"create_data"`
	PageCount      int             `json:"page_count"`
	Width          int             `json:"width"`
	Height         int             `json:"height"`
	SanityLevel    int             `json:"sanity_level"`
	XRestrict      int             `json:"x_restrict"`
	Series         *Series         `json:"series"`
	MetaSinglePage *MetaSinglePage `json:"meta_single_page"`
	MetaPages      []MetaPage      `json:"meta_pages"`
	TotalView      int             `json:"total_view"`
	TotalBookmarks int             `json:"total_bookmarks"`
	IsBookmarked   bool            `json:"is_bookmarked"`
	Visible        bool            `json:"visible"`
	IsMuted        bool            `json:"is_muted"`
	TotalComments  int             `json:"total_comments"`
}

type Series struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type IllustsResponse struct {
	Illusts []Illust `json:"illusts"`
	NextURL string   `json:"next_url"`
}
type IllustResponse struct {
	Illust Illust `json:"illust"`
}

type IllustComments struct {
	TotalComments uint64    `json:"total_comments"`
	Comments      []Comment `json:"comments"`
	NextURL       string    `json:"next_url"`
}

type Comment struct {
	ID             uint64   `json:"id"`
	CommentComment string   `json:"comment"`
	Date           string   `json:"date"`
	User           *User    `json:"user"`
	HasReplies     bool     `json:"has_replies"`
	ParentComment  *Comment `json:"parent_comment"`
}

type IllustCommentAddResult struct {
	Comment Comment `json:"comment"`
}

type IllustRecommended struct {
	Illusts        []Illust       `json:"illusts"`
	RankingIllusts []interface{}  `json:"ranking_illusts"`
	ContestExists  bool           `json:"contest_exists"`
	PrivacyPolicy  *PrivacyPolicy `json:"privacy_policy"`
	NextURL        string         `json:"next_url"`
}

type PrivacyPolicy struct {
}

type TrendingTagsIllust struct {
	TrendTags []TrendTag `json:"trend_tags"`
}

type TrendTag struct {
	Tag            string  `json:"tag"`
	TranslatedName string  `json:"translated_name"`
	Illust         *Illust `json:"illust"`
}

type SearchIllustResult struct {
	Illusts         []Illust `json:"illusts"`
	NextURL         string   `json:"next_url"`
	SearchSpanLimit int      `json:"search_span_limit"`
}

type IllustBookmarkDetail struct {
	BookmarkDetail BookmarkDetail `json:"bookmark_detail"`
}

type BookmarkDetail struct {
	IsBookmarked bool                `json:"is_bookmarked"`
	Tags         []BookmarkDetailTag `json:"tags"`
	Restrict     string              `json:"restrict"`
}

type BookmarkDetailTag struct {
	Name         string `json:"name"`
	IsRegistered bool   `json:"is_registered"`
}

type UserBookmarkTags struct {
	BookmarkTags []interface{} `json:"bookmark_tags"`
	NextURL      string        `json:"next_url"`
}

type UserFollowList struct {
	UserPreviews []UserPreview `json:"user_previews"`
	NextURL      string        `json:"next_url"`
}

type UserPreview struct {
	User    User          `json:"user"`
	Illusts []Illust      `json:"illusts"`
	Novels  []interface{} `json:"novels"`
	IsMuted bool          `json:"is_muted"`
}

type UserList struct {
	Users []interface{} `json:"users"`
}

type UgoiraMetadata struct {
	UgoiraMetadataUgoiraMetadata UgoiraMetadataClass `json:"ugoira_metadata"`
}

type UgoiraMetadataClass struct {
	ZipURLs UserImages `json:"zip_urls"`
	Frames  []Frame    `json:"frames"`
}

type Frame struct {
	File  string `json:"file"`
	Delay int    `json:"delay"`
}

type ShowcaseArticle struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    []Body `json:"body"`
}

type Body struct {
	ID                      string           `json:"id"`
	Lang                    string           `json:"lang"`
	Entry                   Entry            `json:"entry"`
	Tags                    []Tag            `json:"tags"`
	ThumbnailURL            string           `json:"thumbnailUrl"`
	Title                   string           `json:"title"`
	PublishDate             int              `json:"publishDate"`
	Category                string           `json:"category"`
	SubCategory             string           `json:"subCategory"`
	SubCategoryLabel        string           `json:"subCategoryLabel"`
	SubCategoryIntroduction string           `json:"subCategoryIntroduction"`
	Introduction            string           `json:"introduction"`
	Footer                  string           `json:"footer"`
	Illusts                 []BodyIllust     `json:"illusts"`
	RelatedArticles         []RelatedArticle `json:"relatedArticles"`
	FollowingUserIDs        []interface{}    `json:"followingUserIds"`
	IsOnlyOneUser           bool             `json:"isOnlyOneUser"`
}

type Entry struct {
	ID                        string                `json:"id"`
	Title                     string                `json:"title"`
	PureTitle                 string                `json:"pure_title"`
	Catchphrase               string                `json:"catchphrase"`
	Header                    string                `json:"header"`
	Body                      string                `json:"body"`
	Footer                    string                `json:"footer"`
	Sidebar                   string                `json:"sidebar"`
	PublishDate               int                   `json:"publish_date"`
	Language                  string                `json:"language"`
	PixivisionCategorySlug    string                `json:"pixivision_category_slug"`
	PixivisionCategory        PixivisionCategory    `json:"pixivision_category"`
	PixivisionSubcategorySlug string                `json:"pixivision_subcategory_slug"`
	PixivisionSubcategory     PixivisionSubcategory `json:"pixivision_subcategory"`
	Tags                      []Tag                 `json:"tags"`
	ArticleURL                string                `json:"article_url"`
	Intro                     string                `json:"intro"`
	FacebookCount             string                `json:"facebook_count"`
	TwitterCount              string                `json:"twitter_count"`
}

type BodyIllust struct {
	SpotlightArticleID                  int         `json:"spotlight_article_id"`
	IllustID                            int         `json:"illust_id"`
	Description                         string      `json:"description"`
	Language                            string      `json:"language"`
	IllustUserID                        string      `json:"illust_user_id"`
	IllustTitle                         string      `json:"illust_title"`
	IllustExt                           string      `json:"illust_ext"`
	IllustWidth                         string      `json:"illust_width"`
	IllustHeight                        string      `json:"illust_height"`
	IllustRestrict                      string      `json:"illust_restrict"`
	IllustXRestrict                     string      `json:"illust_x_restrict"`
	IllustCreateDate                    string      `json:"illust_create_date"`
	IllustUploadDate                    string      `json:"illust_upload_date"`
	IllustServerID                      string      `json:"illust_server_id"`
	IllustHash                          string      `json:"illust_hash"`
	IllustType                          string      `json:"illust_type"`
	IllustSanityLevel                   int         `json:"illust_sanity_level"`
	IllustBookStyle                     string      `json:"illust_book_style"`
	IllustPageCount                     string      `json:"illust_page_count"`
	IllustCustomThumbnailUploadDatetime interface{} `json:"illust_custom_thumbnail_upload_datetime"`
	IllustComment                       string      `json:"illust_comment"`
	UserAccount                         string      `json:"user_account"`
	UserName                            string      `json:"user_name"`
	UserComment                         string      `json:"user_comment"`
	URL                                 URL         `json:"url"`
	UgoiraMeta                          interface{} `json:"ugoira_meta"`
	UserIcon                            string      `json:"user_icon"`
}

type URL struct {
	The1200X1200    string `json:"1200x1200"`
	The768X1200     string `json:"768x1200"`
	Ugoira600X600   string `json:"ugoira600x600"`
	Ugoira1920X1080 string `json:"ugoira1920x1080"`
}

type RelatedArticle struct {
	ID                        string        `json:"id"`
	Ja                        PrivacyPolicy `json:"ja"`
	En                        PrivacyPolicy `json:"en"`
	Zh                        PrivacyPolicy `json:"zh"`
	ZhTw                      PrivacyPolicy `json:"zh_tw"`
	PublishDate               int           `json:"publish_date"`
	Category                  string        `json:"category"`
	PixivisionCategorySlug    string        `json:"pixivision_category_slug"`
	PixivisionSubcategorySlug string        `json:"pixivision_subcategory_slug"`
	Thumbnail                 string        `json:"thumbnail"`
	ThumbnailIllustID         string        `json:"thumbnail_illust_id"`
	HasBody                   string        `json:"has_body"`
	IsPr                      string        `json:"is_pr"`
	PrClientName              string        `json:"pr_client_name"`
	EditStatus                string        `json:"edit_status"`
	TranslationStatus         string        `json:"translation_status"`
	IsSample                  string        `json:"is_sample"`
	Illusts                   []interface{} `json:"illusts"`
	NovelIDs                  []interface{} `json:"novel_ids"`
	Memo                      string        `json:"memo"`
	FacebookCount             string        `json:"facebook_count"`
	TweetCount                string        `json:"tweet_count"`
	TweetMaxCount             string        `json:"tweet_max_count"`
	Tags                      []interface{} `json:"tags"`
	TagIDs                    interface{}   `json:"tag_ids"`
	NumberedTags              []interface{} `json:"numbered_tags"`
	MainAbtestPatternID       string        `json:"main_abtest_pattern_id"`
	AdvertisementID           string        `json:"advertisement_id"`
}

type PixivisionCategory struct {
	Label        string `json:"label"`
	Introduction string `json:"introduction"`
}

type PixivisionSubcategory struct {
	Label        string `json:"label"`
	LabelEn      string `json:"label_en"`
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
	ImageURL     string `json:"image_url"`
	BigImageURL  string `json:"big_image_url"`
}
