import generic from "../generic";

const API_URL = "/api/policy/";

const Get = () => {
    return generic.GenericGetAll(API_URL)
};

const Update = (policy) => {
    return generic.GenericUpdateNoID(API_URL, policy)
};

export default {
    Get,
    Update
}