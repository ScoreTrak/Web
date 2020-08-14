import axios from "axios";
import authHeader from "../auth/auth-header";

export const getLastNonElapsingRound = () => {
    return axios.get("/api/last_non_elapsing/", {headers: authHeader()}).then( response =>
        {
            if (response.data) {
                return response.data;
            }
        }
    );
};

export default {
    getLastNonElapsingRound
}