import axios from "axios";
import authHeader from "../auth/auth-header";

const API_URL = "/api/report/";

const getReport = () => {
    return axios.get(API_URL, {headers: authHeader()}).then((response) => {
        if (response.data) {
            return response.data;
        }
    });
};

export default {
    getReport
}