# KODE-notes-GoLand
## end-points:

* All CRUD request contains next header. //to-do change to jwt authorization 
Request-header - {
    Authorization: Basic Base64(login:password)
}

If user haven't this header or have some problem - Response code 401

* All unsuccessfull response have error code and message

###POST /api/user/register

Request-body: { 
    login: string
    Password: string
}

If success

Response-body {
    token: string //after changing auth type to jwt
}

### GET api/notes/{sort: {
                           "id" |
                           "-id" |
                           "title" |
                           "-title" |
                           "body" |
                           "-body" |
                           "expir" |
                           "-expire" |
                           "isprivate"|
                           "-isprivate"
                         }
If success
Response-body - {
    {
         id: number 
         title: string
         body: string
         expire: number
         isprivate: bool
    }
}

