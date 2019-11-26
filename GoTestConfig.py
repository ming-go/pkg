import os
import yaml

with open("./GoTestConfig.yml", "r") as f:
    dictConfig = yaml.load(f)

for key in dictConfig:
    os.putenv(key, dictConfig[key])
    print(key, "=", os.getenv(key))
