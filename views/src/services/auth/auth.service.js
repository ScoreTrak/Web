import axios from "axios";

const API_URL = "/auth/";

const login = (username, password) => {
  return axios
    .post(API_URL + "login", {
      username,
      password,
    })
    .then((response) => {
      console.log(response)
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

const getToken = () => {
  return JSON.parse(localStorage.getItem("token"));
};

export default {
  login,
  logout,
  getToken,
};
