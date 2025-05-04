module github.com/chriscow/vapi-go

go 1.18

replace github.com/chriscow/minds => ../thoughtnet/minds

require (
	github.com/chriscow/minds v0.0.7
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/sashabaranov/go-openai v1.39.1 // indirect
)
