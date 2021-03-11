# eviper
Exends viper to rewrite values from another source


## Todos
* Tests
* Lint


```go
package main

import (
	"github.com/spf13/viper"
	"github.com/Alexitdv/eviper"
)

type Opts struct {
	Addr string `env:"ADDR"`
}

func main() {
	c := eviper.New(viper.New())
	c.AutomaticEnv()
	c.AddConfigPath("./configs/local")
	c.SetConfigName("service")
	c.SetConfigType("toml")
	err = c.Unmarshal(&conf) // Unmarshal to struct if needed
}
```