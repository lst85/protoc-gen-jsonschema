package testdata

const ArrayOfPrimitives = `{
    "components": {
        "schemas": {
            "ArrayOfPrimitives": {
                "properties": {
                    "description": {
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "string"
                            }
                        ]
                    },
                    "luckyNumbers": {
                        "items": {
                            "oneOf": [
                                {
                                    "type": "null"
                                },
                                {
                                    "type": "integer"
                                }
                            ]
                        },
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "array"
                            }
                        ]
                    },
                    "luckyBigNumbers": {
                        "items": {
                            "oneOf": [
                                {
                                    "type": "integer"
                                },
                                {
                                    "type": "string"
                                },
                                {
                                    "type": "null"
                                }
                            ]
                        },
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "array"
                            }
                        ]
                    },
                    "keyWords": {
                        "items": {
                            "oneOf": [
                                {
                                    "type": "null"
                                },
                                {
                                    "type": "string"
                                }
                            ]
                        },
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "array"
                            }
                        ]
                    },
                    "bigNumber": {
                        "oneOf": [
                            {
                                "type": "integer"
                            },
                            {
                                "type": "string"
                            },
                            {
                                "type": "null"
                            }
                        ]
                    }
                },
                "oneOf": [
                    {
                        "type": "null"
                    },
                    {
                        "type": "object"
                    }
                ]
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
