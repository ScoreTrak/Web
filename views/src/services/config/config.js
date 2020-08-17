import generic from "../generic";

const API_URL = "/api/config/";

const Get = () => {
    return generic.GenericGetAll(API_URL)
};

const GetStaticConfig = () => {
    return generic.GenericGetAll(API_URL+"static_config")
};

const GetStaticWebConfig = () => {
    return generic.GenericGetAll(API_URL+"static_web_config")
};

const Update = (static_config) => {
    return generic.GenericUpdateNoID(API_URL, static_config)
};


const DeleteCompetition = () => {
    return generic.GenericDeleteNoID(API_URL+"delete_competition")
};

const ResetCompetition = () => {
    return generic.GenericDeleteNoID(API_URL+"reset_competition")
};



export default {
    Get,
    GetStaticConfig,
    GetStaticWebConfig,
    Update,
    ResetCompetition,
    DeleteCompetition
}