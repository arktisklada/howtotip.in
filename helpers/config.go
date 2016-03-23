package helpers

import (
  "bufio"
  "errors"
  "flag"
  "fmt"
  "log"
  "os"
  "strings"
)

type config map[string]string

func readConfigFile(filename string) (map[string]string, error) {
  config := map[string]string{
    "dbhost": "localhost",
  }

  file, err := os.Open(filename)
  if err != nil {
    return nil, errors.New("Couldn't read config " + filename)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := strings.TrimSpace(strings.SplitN(scanner.Text(), "#", 2)[0])

    if line == "" {
      continue
    }

    pieces := strings.SplitN(line, "=", 2)
    if len(pieces) != 2 {
      return nil, errors.New(fmt.Sprintf("Couldn't parse line \"%s\"", line))
    }

    config[strings.TrimSpace(pieces[0])] = strings.TrimSpace(pieces[1])
  }

  return config, nil
}

func ReadConfig() (config) {
  var config_file string
  flag.StringVar(&config_file, "config", "./config.cfg", "config file location")
  flag.Parse()
  cfg, err := readConfigFile(fmt.Sprintf("%v", config_file))
  if err != nil {
    log.Fatal(err.Error())
  }

  if lf := cfg["log_file"]; lf != "" {
    logfile, err := os.OpenFile(lf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0666))
    if err != nil {
      log.Println(err.Error())
    } else {
      log.SetOutput(logfile)
    }
  }

  if pidfile := cfg["pid_file"]; pidfile != "" {
    file, err := os.Create(pidfile)
    if err != nil {
      log.Fatal(err.Error())
    }

    file.WriteString(fmt.Sprintln(os.Getpid()))
    file.Close()
  }

  return cfg
}
