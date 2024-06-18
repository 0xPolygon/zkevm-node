package config

// DefaultValues is the default configuration
const DefaultValues = `
[Online]
URI = "localhost:6900"
StreamType = 1

[Offline]
Port = 6901
Filename = "datastream.bin"
Version = 3
ChainID = 1440
WriteTimeout = "5s"
UpgradeEtrogBatchNumber = 0

[StateDB]
User = "state_user"
Password = "state_password"
Name = "state_db"
Host = "localhost"
Port = "5432"
EnableLog = false	
MaxConns = 200

[MerkleTree]
URI = "localhost:50061"
MaxThreads = 0
CacheFile = ""

[Log]
Environment = "development" # "production" or "development"
Level = "info"
Outputs = ["stderr"]
`
