# Reverse engineering of the site lkdr.nalog.ru

Here will be a dump of my findings that I learned while looking at what requests the site sends to the client and back.

## How the site authorizes the user and what it does after authorization

### Short version

1. The user enters his phone number in the input form
2. Next, the user receives an SMS, which he must enter in the form
3. If the user enters the code correctly, he is given an authorization token 

Short version of XHR requests: 

```
GET profile -> GET token -> POST start <-> POST verify -> GET identifiers -> GET profile -> GET receipt
```

### Long version

---

When a client first visits the site, there is an XHR request right away:

- Url                  : `https://lkdr.nalog.ru/api/v1/user/profile`
- Type                 : `GET`
- Expected status code : `401`

##### Server's response: 

```json
{
   "code":"Not Authenticated",
   "message":"Not Authenticated",
   "additionalInfo":{
      
   }
}
```

Immediately following with the next request, but this time with the payload:

- URL                 : `https://lkdr.nalog.ru/api/v1/user/profile`
- Type                : `GET`
- Expected status code: `422`

##### Client's payload:

```json
{
   "refreshToken":null,
   "deviceInfo":{
      "sourceDeviceId":"device_id",
      "sourceType":"WEB",
      "appVersion":"1.0.0",
      "metaDetails":{
         "userAgent":"browser"
      }
   }
}
```

##### Server's response: 

```json
{
   "code":"validation.failed",
   "message":"Пустой refreshToken",
   "additionalInfo":null
}
```

---

When we enter the phone numner into the input form, we start the `SMS Challenge`:

- URL : `https://lkdr.nalog.ru/api/v1/auth/challenge/sms/start`
- Type: `POST`

##### Client's payload:

```json
{"phone":"123"}
```

##### Server's response:

```json
{
   "challengeToken":"UUID",
   "challengeTokenExpiresIn":"2021-01-01T00:00:00.000Z",
   "challengeTokenExpiresInSec":120000
}
```

---

Then we're sending our SMS code with `JSON` info, that we received earlier:

- URL : `https://lkdr.nalog.ru/api/v1/auth/challenge/sms/verify`
- Type: `POST`

##### Client's payload

```json
{
   "challengeToken":"UUID",
   "phone":"123",
   "code":"6 digits code",
   "deviceInfo":{
      "sourceDeviceId":"deviceid",
      "sourceType":"WEB",
      "appVersion":"1.0.0",
      "metaDetails":{
         "userAgent":"useragent"
      }
   }
}
```

##### Server's response

```json
{
   "refreshToken":"long ass refresh token",
   "refreshTokenExpiresIn":null,
   "token":"long ass token",
   "tokenExpireIn":"2021-01-01T00:00:00.000Z",
   "profile":{
      "taxpayerPerson":{
         "email":null,
         "phone":"123",
         "inn":null,
         "fullName":null,
         "shortName":null,
         "status":"BASIC_ACTIVE",
         "address":null,
         "oktmo":null,
         "authorityCode":null,
         "firstName":null,
         "lastName":null,
         "middleName":null
      },
      "authType":"SMS"
   }
}
```

---

The service then automatically makes a GET request

- URL : `https://lkdr.nalog.ru/api/v1/identifiers`
- Type: `GET`

##### Server's response

```json
{
   "identifiers":[
      {
         "login":"123",
         "type":"SMS"
      }
   ],
   "pendingConfirmation":[
      {
         "login":"some@email.com",
         "type":"EMAIL",
         "confirmationExpiresAt":"2021-01-01T00:00:00.000Z",
      },
      {
         "login":"some_other_email@email.com",
         "type":"EMAIL",
         "confirmationExpiresAt":"2021-01-01T00:00:00.000Z",
      }
   ]
}
```

---

Then gets information about the user's profile:

- URL : `https://lkdr.nalog.ru/api/v1/user/profile`
- Type: `GET`

##### Server's response: 

```json
{
   "user":{
      "taxpayerPerson":{
         "email":null,
         "phone":"123",
         "inn":null,
         "fullName":null,
         "shortName":null,
         "status":"BASIC_ACTIVE",
         "address":null,
         "oktmo":null,
         "authorityCode":null,
         "firstName":null,
         "lastName":null,
         "middleName":null
      },
      "authType":"SMS"
   }
}
```

---

Then the site makes a request for receipts that belong to the user:

- URL : `https://lkdr.nalog.ru/api/v1/receipt`
- Type: `GET`

##### Client's payload:

```json
{
   "limit":10,
   "offset":0,
   "dateFrom":null,
   "dateTo":null,
   "orderBy":"RECEIVE_DATE:DESC",
   "inn":null
}
```

##### Server's response: 

```json
{
   "receipts":[
      {
         "buyer":"123",
         "buyerType":"1",
         "createdDate":"2021-01-01T00:00:00",
         "fiscalDocumentNumber":"123456",
         "fiscalDriveNumber":"16-digit number",
         "fiscalSign":null,
         "totalSum":"5.00",
         "kktOwner":"ООО \"Пятерочка-крипто\"",
         "kktOwnerInn":"INN Number",
         "key":"123|1|2021 01 01 00:00:00|123456|16-digit number"
      },
      {
         "buyer":"123",
         "buyerType":"1",
         "createdDate":"2021-01-01T00:00:00",
         "fiscalDocumentNumber":"12345",
         "fiscalDriveNumber":"16-digit number",
         "fiscalSign":null,
         "totalSum":"69000.00",
         "kktOwner":"ООО \"Смерть\"",
         "kktOwnerInn":"INN Number",
         "key":"123|1|2021 01 01 00:00:00|12345|16-digit number"
      }
   ],
   "hasMore":true
}
```

---

And then it makes a request every 5000 milliseconds* to get some kind of notification.

* Obtained from `const timer = setInterval(checkNewNotifications, 5000);`

- URL : `https://lkdr.nalog.ru/api/v1/notification/count`
- Type: `GET`

##### Server's response

```json
{"numberUnacknowledgedNotifications":0}
```
