import axios from "axios";
import jwt_decode from "jwt-decode"

const API_URL = "/auth/";

const login = (username, password) => {
  return axios
    .post(API_URL + "login", {
      username,
      password,
    })
    .then((response) => {
      if (response.data.token) {
        localStorage.setItem("token", JSON.stringify(response.data.token));
        return response.data.token;
      } else {
        throw new Error("We were unable to authenticate this account");
      }
    });
};



const logout = () => {
  localStorage.removeItem("token");
};

const logAnonymousUser = () => {
  return login("", "");
};

const logAnonymousUserIfNecessary = () => {
  if (!isAValidToken()){
    return logAnonymousUser("", "")
  }
}

const getToken = () => {
  return JSON.parse(localStorage.getItem("token"));
};

const tokenExists = () =>{
  return !!localStorage.getItem("token")
}

const getDecodedJWT = () => {
  try{
    return jwt_decode(localStorage.getItem("token"))
  }
  catch (err){
    return {}
  }
}

const getCurrentRole = () => {
  return getDecodedJWT()["role"]
};

const getCurrentTeamID = () => {
  return getDecodedJWT()["team_id"]
};

const isAValidRole = () => {
  const role = getCurrentRole()
  return (role && role !== "anonymous")
}

const tokenExpired = () => {
  let current_time = new Date().getTime() / 1000;
  return current_time > getDecodedJWT().exp
}

const isAValidToken = () => {
  return (tokenExists() && !tokenExpired())
}


export default {
  logAnonymousUserIfNecessary,
  logAnonymousUser,
  login,
  logout,
  getToken,
  getCurrentRole,
  getDecodedJWT,
  isAValidRole,
  isAValidToken,
  getCurrentTeamID
};
