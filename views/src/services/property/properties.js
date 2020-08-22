import generic from "../generic"

const API_URL = "/api/property/";
const API_URL_MULTI = "/api/properties/";



const Update = (service_id, key, property) => { return generic.GenericUpdate(API_URL, `${service_id}/${key}`, property) };
const GetAllByServiceIDKey = (service_id, key,) => { return generic.GenericGetByID(API_URL, `${service_id}/${key}`) };
const Delete = (service_id, key) => { return generic.GenericDelete(API_URL, `${service_id}/${key}`) };
const Create = (property) => { return generic.GenericCreate(API_URL, property) };

const GetAllByServiceID = (service_id) => { return generic.GenericGetByID(API_URL_MULTI, service_id) };
const GetAll = () => { return generic.GenericGetAll(API_URL_MULTI) };



export default {
    Create,
    Update,
    GetAllByServiceIDKey,
    GetAll,
    Delete,
    GetAllByServiceID
};
