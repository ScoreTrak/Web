import axios from "axios";

const API_URL = "/api/user";

const createUser = (username, password) => {
    return axios.post(API_URL, [{
            username,
            password,
        }]);
};

export default {
    createUser
};
