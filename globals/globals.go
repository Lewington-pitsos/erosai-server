package globals

const BetServerPort = ":8082"
const NsfwCutoff = 70
const MLServerEndpoint = "http://localhost:8001"

const ImageSaveDirectory = "./temp/"

// Additional globals needed by other packages
var LogLevel = "info"
var DefaultAccessToken = "default_token"
var NZVenues = []string{}
var Password = "default_password"

type BookieName string