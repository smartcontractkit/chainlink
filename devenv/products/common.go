package products

// import (
// 	"os"
// 	"strings"
// )

// // Load loads TOML configurations from environment variable, ex.: CTF_CONFIGS=env.toml,overrides.toml
// // and unmarshalls the files from left to right overriding keys.
// func Load[T any](path string) (*T, error) {
// 	var config T
// 	paths := strings.Split(os.Getenv(EnvVarTestConfigs), ",")
// 	for _, path := range paths {
// 		data, err := os.ReadFile(path) 
// 		if err != nil {
// 			retu
// 		}

// 		decoder := toml.NewDecoder(strings.NewReader(string(data)))
// 		decoder.DisallowUnknownFields()

// 		if err := decoder.Decode(&config); err != nil {
// 			var details *toml.StrictMissingError
// 			if errors.As(err, &details) {
// 				fmt.Println(details.String()) //nolint:forbidigo
// 			}
// 			return nil, fmt.Errorf("failed to decode TOML config, strict mode: %s", err)
// 		}
// 	}
// 	if L.GetLevel() == zerolog.TraceLevel {
// 		L.Trace().Msg("Merged inputs")
// 		spew.Dump(config) //nolint:forbidigo
// 	}
// 	return &config, nil
// }