import generic from "../generic"

const API_URL = "/api/host_group/";


const Create = (host) => { return generic.GenericCreate(API_URL, host) };
const Update = (id, host) => { return generic.GenericUpdate(API_URL, id, host) };
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
