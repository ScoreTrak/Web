import generic from "../generic";

const API_URL = "/api/report/";

const Get = () => {
    return generic.GenericGetAll(API_URL)
};

export default {
    Get
}