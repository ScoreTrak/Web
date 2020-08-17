import generic from "../generic"

const API_URL = "/api/user/";


const Create = (users) => { return generic.GenericCreate(API_URL, users) };
const Update = (id, user) => { return generic.GenericUpdate(API_URL, id, user) };
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
