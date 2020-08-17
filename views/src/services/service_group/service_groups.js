import generic from "../generic"

const API_URL = "/api/service_group/";


const Create = (service_group) => { return generic.GenericCreate(API_URL, service_group) };
const Update = (id, service_group) => { return generic.GenericUpdate(API_URL, id, service_group) };
const GetByID = (id) => { return generic.GenericGetByID(API_URL, id) };
const GetAll = () => { return generic.GenericGetAll(API_URL) };
const Delete = (id) => { return generic.GenericDelete(API_URL, id) };


export default {
    Create,
    Update,
    GetByID,
    GetAll,
    Delete
};
