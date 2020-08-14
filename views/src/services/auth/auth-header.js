import AuthService from "./auth"

export default function authHeader() {
  const token = AuthService.getToken()
  if (token) {
    return { Authorization: 'Bearer ' + token };
     } else {
    return {};
  }
}
