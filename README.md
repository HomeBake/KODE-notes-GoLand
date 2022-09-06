# KODE-notes-GoLand
## end-points:

* All CRUD request contains next header. //to-do change to jwt authorization 
Request-header - {
    Authorization: Basic Base64(login:password)
}

If user haven't this header or have some problem - Response code 401

* All unsuccessfull response have error code and message

 Endpoints
  
    GET API/notes           RESPONSE-body -   
                                            {notes:  
    											{  
												id: number,  
												body: string,  
												title: string,  
												isprivate: boolean,  
												expire: int,  
												userId: int,  
										    },...  
										}   
    GET API/notes/{noteId}  RESPONSE-body - {  
                                                id: number,  
												body: string,  
												title: string,  
												isprivate: boolean,  
												expire: int,  
												userId: int,  
											}  
    
    POST API/notes          REQUEST-body - {  
                                                body: string,  
												title: string,  
												isprivate: boolean,  
												expire: int,  
										    }  
				            RESPONSE-body: {result: boolean}  
    PUT API/notes           REQUEST-body - {  
												id: number,  
												body: string,
                                                title: string,  
												isprivate: boolean,  
												expire: int,  
										    }  
						    RESPONSE-body: {result:boolean}  
    DELETE API/notes/{noteId}   
                            RESPONSE-body - {result: boolean}  
    POST API/register       REQUEST-body - {  
											login: string,  
											password: string,  
	                                        }  
						    RESPONSE-body: {result: boolean}  
    POST API/login           
							REQUEST-body - {  
											login: string,  
											password: string, 
											}  
							RESPONSE-body:  {  
											id: int,   
											login: string, 
											password: string,  
											}  
    POST API/notes/{id}     REQUEST-body: {  
											USERACCESSID: int,  
											MODE: string,  
											}  
                            RESPONSE-body: {result: boolean}   

