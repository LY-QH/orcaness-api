config_file=config/app.prod.yaml

if [ ! -f $config_file ]; then
  echo "The configuration file [$config_file] does not exist, please add it first."
  exit 1
fi

GIN_MODE=release go build -o bin/orcaness-api

[ "$?" = "0" ] && echo "The executable file [bin/orcaness-api] has built."