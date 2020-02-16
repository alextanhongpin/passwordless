# passwordless

Simple implementation of passwordless with golang.

It starts with a simple problem statement:

```
As a user,
I want to login without remembering my password,
In order to make login more seamless.
```

## UseCase
### User Request Code

1. User navigates to the passwordless login page.
2. System requests an identifier, email or phone number.
	- Business rule: User can choose either email or mobile for the code to be sent to.
3. User submits identifier.
4. System verifies that the identifier exists, is confirmed and sends the code to the user.
5. User receives the code.

3a. User submits multiple item.
	- System throttles the user requests.

### User Submits Code 

1. User navigates to the passwordless login page with code.
2. User enters the code that was received earlier.
3. System verifies that the code has not yet been used, and has not yet expired.
4. System returns auth tokens to user.
5. System stores the token securely locally and redirect user to the protected pages.

2a. User submits invalid code.
	- System returns error message indicating failure.
	- System throttles the user when there are too many failed requests.

## API

### Passwordless Start

Three options are available:
- send a verification code using email
- send a link using email
- send a verification code using SMS 

Request Parameters:
| Parameter | Description |
| - | - |
| `connection` | How to send the code/link to the user. Can be `email` or `sms` |
| `email` | Set this when `connection=email` |
| `phone_number` | Set this when `connection=sms` |
| `send` | Use `link` to send a link or `code` to send a verification code. If null, a link will be send. |

```
POST /passwordless/start
Content-Type: application/json
{
  "connection": "email|sms",
  "email": "EMAIL",
  "phone_number": "PHONE_NUMBER",
  "send": "link|code"
}
```

### Passwordless Verify

Perform verification of the code send to the user, with `username=email|phone_number` and `otp=verification_code`. This endpoint should be rate-limited to 50 requests per hour per IP to prevent abuse.

```
POST /passwordless/verify
Content-Type: application/json
{
  "realm": "email|sms",
  "grant_type": "passwordless_otp",
  "username": "email|phone_number",
  "otp": "verification_code",
  "send": "link|code"
}
```
