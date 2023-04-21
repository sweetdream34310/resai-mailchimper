# Project: Awaymail
# 📁 Collection: aways


## End-point: create user away
### Method: POST
>```
>{{dev}}/v2/away
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "title": "vicky inbox",
    "activate_allow": "2022-01-22T20:05:16.167Z",
    "repeat": [],
    "is_enabled": true,
    "all_day": true,
    "allowed_users": [],
    "allowed_subjects": []
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get list user away
### Method: GET
>```
>{{dev}}/v2/aways
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get user away by ID
### Method: GET
>```
>{{dev}}/v2/away/61ec766ff5922d0304d8f28a
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: delete user away by ID
### Method: DELETE
>```
>{{dev}}/v2/away/17dd424bb3f823b3
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: update user away
### Method: PATCH
>```
>{{dev}}/v2/away
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "_id": "61bf84ddf5922d90127cc43a",
    "title": "test",
    "activate_allow": "2021-12-20T21:58:09.413Z",
    "deactivate_allow": "2021-12-25T21:58:09.413Z",
    "start_time": "2021-12-20T21:58:09.413Z",
    "end_time": "2021-12-25T21:58:09.413Z",
    "repeat": [],
    "all_day": false,
    "allowed_users": [
        {
            "id": "61b31895f5922d684cf3cd59",
            "name": "casbu",
            "activated": false
        }
    ],
    "allowed_subjects": [
        {
            "id": "61b31895f5922d684cf3cd63",
            "name": "Founders Circle",
            "activated": false
        },
        {
            "id": "61b31895f5922d684cf3cd60",
            "name": "Feedback",
            "activated": false
        }
    ]
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: Inbox


## End-point: get User Inbox
### Method: GET
>```
>{{dev}}/v2/inbox
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body formdata

|Param|value|Type|
|---|---|---|
|skip|0|text|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get User Sent
### Method: GET
>```
>{{dev}}/v2/sent
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body formdata

|Param|value|Type|
|---|---|---|
|skip|0|text|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get User email by message_id
### Method: GET
>```
>{{dev}}/v2/message/17e7f7bb8451eb86
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: send email
### Method: POST
>```
>{{dev}}/v2/send
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "to": "vicky@cloudsource.io",
    "message": "Hello wolrd!",
    "subject": "test send email 2",
    "attachments_url": []
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: update read inbox
### Method: PATCH
>```
>{{dev}}/v2/message/17dd424bb3f823b3
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "is_read": true
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: file


## End-point: upload file
### Method: POST
>```
>{{dev}}/v2/file
>```
### Body formdata

|Param|value|Type|
|---|---|---|
|file|/Users/viduka/Desktop/hello.rtf|file|


### 🔑 Authentication noauth

|Param|value|Type|
|---|---|---|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: user


## End-point: Get Session
### Method: POST
>```
>{{dev}}/v2/session
>```
### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Body (**raw**)

```json
{
    "provider": "gmail",
    "access_token": "ya29.A0ARrdaM8XBy1Mo0JvXQRJcwy5eOX0CUly6n0wHhtdPD7za5AgqeGt7f2FydcRBau-ud93fNlFLc89gmY1KFo_ndmy-eBujq8axp5oX5ray7NhDMNGIRdHN7DYrxb8nyamN2QbxHX6vsodb8Met6Kmwaa-MZIb8g",
    "refresh_token": "1//065YlVvrYOM7VCgYIARAAGAYSNwF-L9IrJX1JhcrAwwS7XFv6qXUSCuc9CdyB5G0UQ4qyVNLfaG-vVpeInGhkg3wxdg22FNx0LMs"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: contacts


## End-point: add contacts
### Method: POST
>```
>{{dev}}/v2/contacts
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "name": "",
    "email": ""
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get contacts list
### Method: GET
>```
>{{dev}}/v2/contacts
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: update contact
### Method: PUT
>```
>{{dev}}/v2/contacts/61fc3d75f5922da337427559
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "name": "jaya",
    "email": "jaya@cloudsource.io"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: remove contact
### Method: DELETE
>```
>{{dev}}/v2/contacts/61fc3d75f5922da337427559
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: userV3


## End-point: Get Session
### Method: POST
>```
>{{dev}}/v3/session
>```
### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Body (**raw**)

```json
{
    "provider": "gmail",
    "swift_token": "486f170b6970e012d8efe5f1871a27ce190da3555a4ce46e64da49411fe65ae8",
    "access_token": "ya29.A0ARrdaM_g5FxoyS0z9HrHE4dvvI5qZcG0AWCoiveNZvYohVkn_lzniO5vI7UGmUoA6B83KVNZ0s_dYEDeTBlTCBxy3SjAPxd6NW5RAM1ygUg3IeaWjvgcRxoDpBlJ6kBZIucSHcO_mCzNGHla8yDRPt7eSTVk",
    "refresh_token": "1//0gtbouq0hVuHoCgYIARAAGBASNwF-L9IrJpXVBPLqCBLm-6vRDrh-bSKr2LgsKafo1aMXarKo2QjlwuUKviaxRjEzzm5-Km08VB4",
    "dev_token_key": false
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: logout
### Method: GET
>```
>{{dev}}/v3/logout
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get profile
### Method: GET
>```
>{{dev}}/v3/user/profile
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: awaysV3


## End-point: create user away
### Method: POST
>```
>{{dev}}/v3/away
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|website|


### Body (**raw**)

```json
{
    "title": "all day away",
    "repeat": ["Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"],
    "activate_allow": "2022-02-10T06:00:00Z",
    "deactivate_allow": "2022-02-10T23:00:00Z",
    "is_enabled": true,
    "all_day": true,
    "message": "Thanks for your message",
    "allowed_contacts": ["xbleder@gmail.com"],
    "allowed_keywords": ["bitrise", "\"away\""]
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get list user away
### Method: GET
>```
>{{dev}}/v3/aways
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get away by ID
### Method: GET
>```
>{{dev}}/v3/away/61fe9023f5922d126a515f11
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: update away
### Method: PUT
>```
>{{dev}}/v3/away
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "id": "6310f974f5922d71e9d13e39",
    "title": "all day away",
    "activate_allow": "2022-02-10T13:00:00+07:00",
    "deactivate_allow": "2022-02-11T06:00:00+07:00",
    "repeat": [
        "Sunday",
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday"
    ],
    "is_enabled": false,
    "all_day": true,
    "message": "Thanks for your message",
    "allowed_contacts": [
        "xbleder@gmail.com"
    ],
    "allowed_keywords": ["bitrise", "\"away\""]
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: enable away
### Method: PATCH
>```
>{{dev}}/v3/away/61fe9023f5922d126a515f11
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "is_enabled": true
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: delete away by ID
### Method: DELETE
>```
>{{dev}}/v3/away/61fe90c5f5922d1a5f31a969
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: enable all aways
### Method: PATCH
>```
>{{dev}}/v3/away
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "is_enabled": true
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: contactsV3


## End-point: get contacts list
### Method: GET
>```
>{{dev}}/v3/contacts
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: update contact
### Method: PUT
>```
>{{dev}}/v3/contacts/61fd1c5cf5922d758a7b74e9
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "name": "jaya k",
    "email": "jaya@cloudsource.io"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: add contacts
### Method: POST
>```
>{{dev}}/v3/contacts
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "name": "",
    "email": ""
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: remove contact
### Method: DELETE
>```
>{{dev}}/v3/contacts/61fc3d75f5922da337427559
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get recent contacts list
### Method: GET
>```
>{{dev}}/v3/contacts/recent
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: inboxV3


## End-point: get User Inbox
### Method: GET
>```
>{{dev}}/v3/inbox?page=1&search=
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Query Params

|Param|value|
|---|---|
|page|1|
|search||


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get User Sent
### Method: GET
>```
>{{dev}}/v3/sent
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "page": 1,
    "search": ""
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get inbox message by message_id
### Method: GET
>```
>{{dev}}/v3/message/17e7f7bb8451eb86
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get sent box message by message_id
### Method: GET
>```
>{{dev}}/v3/sent/message/17e7f7bb8451eb86
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: send email
### Method: POST
>```
>{{dev}}/v3/send
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "to": "xbleder@gmail.com",
    "message": "Good!",
    "subject": "test 2nd multiple email",
    "attachments_url": []
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: send reply email
### Method: POST
>```
>{{dev}}/v3/send
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "to": "xbleder@gmail.com",
    "message": "Good!\r\n\r\nOn Sat, 26 Feb 2022 at 07:57, <vicky@cloudsource.io> wrote:\r\n\r\n> It's working!\r\n\r\nOn Sat, 26 Feb 2022 at 07:57, <vicky@cloudsource.io> wrote:\r\n\r\n> nice one too\r\n>\r\n",
    "subject": "Re: test 2nd multiple email",
    "attachments_url": [],
    "thread_id": "17f33872e81f0f97",
    "message_id": "<CAEhbSB+Go2n4=0bsSqz4LieL1NmMY6cSxoVSprODr7j76dzoOQ@mail.gmail.com>"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: update read inbox
### Method: PATCH
>```
>{{dev}}/v3/message/17dd424bb3f823b3
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "is_read": true
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: delete message inbox
### Method: DELETE
>```
>{{dev}}/v3/message/17ecb5292cc5a7aa
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "is_read": true
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get inbox label
### Method: GET
>```
>{{dev}}/v3/label/INBOX
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Archive Inbox
### Method: POST
>```
>{{dev}}/v3/archive
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: delete archive Inbox
### Method: DELETE
>```
>{{dev}}/v3/archive/inbox
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Archive Inbox by Label
### Method: POST
>```
>{{dev}}/v3/inbox/archive
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "message_ids": ["1801f367cd77a502"]
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get Archive Label count
### Method: GET
>```
>{{dev}}/v3/label/AwayMailArchive
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: direct Google API


## End-point: Get inbox
### Method: GET
>```
>{{dev}}/v3/user/message/inbox?nextPageToken=&prevPageToken=&search=
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Query Params

|Param|value|
|---|---|
|nextPageToken||
|prevPageToken||
|search||


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Get draft
### Method: GET
>```
>{{dev}}/v3/user/message/draft?nextPageToken=&prevPageToken=&search=
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Query Params

|Param|value|
|---|---|
|nextPageToken||
|prevPageToken||
|search||


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Get sent
### Method: GET
>```
>{{dev}}/v3/user/message/sent?nextPageToken=&prevPageToken=&search=
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Query Params

|Param|value|
|---|---|
|nextPageToken||
|prevPageToken||
|search||


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Get trash
### Method: GET
>```
>{{dev}}/v3/user/message/trash?nextPageToken=&prevPageToken=&search=
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Query Params

|Param|value|
|---|---|
|nextPageToken||
|prevPageToken||
|search||


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Get spam
### Method: GET
>```
>{{dev}}/v3/user/message/spam?nextPageToken=&prevPageToken=&search=
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Query Params

|Param|value|
|---|---|
|nextPageToken||
|prevPageToken||
|search||


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get inbox count
### Method: GET
>```
>{{dev}}/v3/label/inbox
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get sent count
### Method: GET
>```
>{{dev}}/v3/label/sent
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get spam count
### Method: GET
>```
>{{dev}}/v3/label/spam
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get draft count
### Method: GET
>```
>{{dev}}/v3/label/draft
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get trash count
### Method: GET
>```
>{{dev}}/v3/label/trash
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: tag message
### Method: POST
>```
>{{dev}}/v3/user/message/tag
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "message_ids": ["17f226eead5320b0", "17f10eca7a3d4f7d"],
    "label_ids": ["Label_1", "Label_7858365802877513870"]
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get attachment data
### Method: GET
>```
>{{dev}}/v3/message/attachment?filename=tree2.jpeg&mime_type=image/jpeg&message_id=17f37069e60762c9&attachment_id=ANGjdJ-6bsDXrXFHuEDJMtTefEmBDQCG6fDiNuqPeI-5kcnaHi3M5gos8voJJnoWmAl1z6qIwZru_-sFKLCZg5mPJEKdM9x1o1kPfj2iGfw3nSckmq6e4Q3SMD8PeINwP5OVcpl8k0i5gFnxZt-pWm9a7c6tQ-vuztfAstmBZYNla3lrcXvMWWtTU4KiFm9ibx9csJfB0oskCDuFhydq5ud1Jd6yMgr1OJXrbYzgS4rKOc9j_faSbaQq-7oZqSglJbv_IsqAbWu3ypq9PRAOvmJGZkgdEpyTJN9Odn4Wdin8e-L4ZpC0hiN_I8z0R4BPSdl8lMrnj_Dd2UNimGvQEbdxshZXBqQpKHF8jY_AZlC98LceM5HuSNTfrNWlnTot33oayjABTduL2wwWRKTn
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json

```

### Query Params

|Param|value|
|---|---|
|filename|tree2.jpeg|
|mime_type|image/jpeg|
|message_id|17f37069e60762c9|
|attachment_id|ANGjdJ-6bsDXrXFHuEDJMtTefEmBDQCG6fDiNuqPeI-5kcnaHi3M5gos8voJJnoWmAl1z6qIwZru_-sFKLCZg5mPJEKdM9x1o1kPfj2iGfw3nSckmq6e4Q3SMD8PeINwP5OVcpl8k0i5gFnxZt-pWm9a7c6tQ-vuztfAstmBZYNla3lrcXvMWWtTU4KiFm9ibx9csJfB0oskCDuFhydq5ud1Jd6yMgr1OJXrbYzgS4rKOc9j_faSbaQq-7oZqSglJbv_IsqAbWu3ypq9PRAOvmJGZkgdEpyTJN9Odn4Wdin8e-L4ZpC0hiN_I8z0R4BPSdl8lMrnj_Dd2UNimGvQEbdxshZXBqQpKHF8jY_AZlC98LceM5HuSNTfrNWlnTot33oayjABTduL2wwWRKTn|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Get Detail Inbox
### Method: GET
>```
>{{dev}}/user/thread/message/:id
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Get Archive Label
### Method: GET
>```
>{{dev}}/v3/user/message/archive?nextPageToken=&prevPageToken=&search=
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Query Params

|Param|value|
|---|---|
|nextPageToken||
|prevPageToken||
|search||


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get Archive Label count
### Method: GET
>```
>{{dev}}/v3/label/archive
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Collection: Label


## End-point: label list
### Method: GET
>```
>{{dev}}/v3/label
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: create label
### Method: POST
>```
>{{dev}}/v3/label
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
    "name": "testLabel"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: patch label
### Method: PATCH
>```
>{{dev}}/v3/label/:labelID
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### Body (**raw**)

```json
{
  "name": "testPatch"
}
```

### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: delete label
### Method: DELETE
>```
>{{dev}}/v3/label/:labelID
>```
### Headers

|Content-Type|Value|
|---|---|
|X-Client|gmail|


### Headers

|Content-Type|Value|
|---|---|
|X-agent|ios|


### 🔑 Authentication bearer

|Param|value|Type|
|---|---|---|
|token|{{token}}|string|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
_________________________________________________
Powered By: [postman-to-markdown](https://github.com/bautistaj/postman-to-markdown/)
