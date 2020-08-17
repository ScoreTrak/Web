import axios from "axios";
import authHeader from "./auth/auth-header";

const GenericCreate = (API_URL,data) => {
    return axios.post(API_URL, data, {headers: authHeader()})
};

const GenericUpdate = (API_URL,id, data) => {
    return axios.patch(API_URL+id, data,{headers: authHeader()});
};

const GenericUpdateNoID = (API_URL, data) => {
    return axios.patch(API_URL, data,{headers: authHeader()});
};

const GenericDownloadFile = (API_URL, filename) => {
    return axios.get(API_URL, {headers: authHeader(), responseType: 'blob'}).then( response =>
    {   const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', filename);
        document.body.appendChild(link);
        link.click();
        link.remove()
    });
};

const GenericGetByID = (API_URL,id) => {
    return axios.get(API_URL+id, {headers: authHeader()}).then( response =>
    {
        if (response.data == null){
             return {}
        }
        return response.data
    });
};

const GenericGetAll = (API_URL) => {
    return axios.get(API_URL, {headers: authHeader()}).then( response =>
    {
        if (response.data == null){
            return []
        }
        return response.data

    })
};

const GenericDelete = (API_URL,id) => {
    return axios.delete(API_URL+id,{headers: authHeader()});
};

const GenericDeleteNoID = (API_URL) => {
    return axios.delete(API_URL,{headers: authHeader()});
};

const GenericUploadFile = (API_URL,formData) => {
    return axios.post(API_URL,formData,{headers: {...authHeader(), 'Content-Type': 'multipart/form-data'}});
};


export default {
    GenericDeleteNoID,
    GenericCreate,
    GenericUpdate,
    GenericGetByID,
    GenericGetAll,
    GenericDelete,
    GenericUpdateNoID,
    GenericUploadFile,
    GenericDownloadFile

};
