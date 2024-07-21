wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{"content": "test", "status_id": 1}'

function request()
    return wrk.format(nil, "/api/message")
end
