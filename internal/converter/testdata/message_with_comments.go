package testdata

const MessageWithComments = `{
    "components": {
        "schemas": {
            "MessageWithComments": {
                "properties": {
                    "name1": {
                        "type": "string",
                        "description": "This field is supposed to represent blahblahblah"
                    }
                },
                "type": "object",
                "description": "This is a message level comment and talks about what this message is and why you should care about it!"
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
