package main

import "net/http"

var requestType = remoteRequest

const defaultPort = "2040"
const defaultMethod = http.MethodPost

const remoteRequest = "REMOTE"
const localRequest = "LOCAL"

const defaulthpserver = "http://127.0.0.1:2040"
const httpreserveFunction = "/httpreserve"
