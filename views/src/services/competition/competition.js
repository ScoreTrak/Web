import generic from "../generic";

const API_URL = "/api/competition/";

const FetchCoreCompetition = () => {
    return generic.GenericDownloadFile(API_URL+"export_core", 'competition_core.yaml')
};

const FetchEntireCompetition = () => {
    return generic.GenericDownloadFile(API_URL+"export_all", 'competition_full.yaml')
};

const LoadCompetition = (formData) => {
    return generic.GenericUploadFile(API_URL+"upload", formData)
};

export default {
    LoadCompetition,
    FetchEntireCompetition,
    FetchCoreCompetition
}