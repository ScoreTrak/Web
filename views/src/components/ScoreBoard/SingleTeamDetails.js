import React, { useState} from "react";
import Accordion from "@material-ui/core/Accordion";
import AccordionSummary from "@material-ui/core/AccordionSummary";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import Typography from "@material-ui/core/Typography";
import AccordionDetails from "@material-ui/core/AccordionDetails";
import Box from "@material-ui/core/Box";
import {makeStyles} from "@material-ui/core/styles";
import CheckCircleOutlineIcon from '@material-ui/icons/CheckCircleOutline';
import ErrorIcon from '@material-ui/icons/Error';
import Alert from "@material-ui/lab/Alert";
import AlertTitle from "@material-ui/lab/AlertTitle";
import Grid from "@material-ui/core/Grid";
import MaterialTable from "material-table";
import PropertyService from "../../services/property/properties";
import properties from "../../services/property/properties";


const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        width: '100%',

    },

    customAccordionSuccessHeader: {
        borderRight: `1px solid ${theme.palette.success.main}`,
        borderLeft: `1px solid ${theme.palette.success.main}`,
    },
    customAccordionErrorHeader: {
        borderRight: `1px solid ${theme.palette.error.main}`,
        borderLeft: `1px solid ${theme.palette.error.main}`,
    },

    paper: {
        padding: theme.spacing(1),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },

    button: {
        margin: theme.spacing(1),
    },

    input: {
        display: 'none',
    },

    heading: {
        fontSize: theme.typography.pxToRem(15),
        flexBasis: '33.33%',
        flexShrink: 0,
    },
    secondaryHeading: {
        fontSize: theme.typography.pxToRem(15),
        color: theme.palette.text.secondary,
    },
    iconSuccess: {
        color: theme.palette.success.main,
        marginRight: 12,
        fontSize: 22,
        opacity: 0.9,
    },

    iconError: {
        color: theme.palette.error.main,
        marginRight: 12,
        fontSize: 22,
        opacity: 0.9,
    },

}));

export default function SingleTeamDetails(props) {
    const classes = useStyles();
    async function reload(service) {
        const results = await PropertyService.GetAllByServiceID(service)
        if ((!!results) && (results.constructor === Object)){
            return []
        }
        return [...results]
    }
    const columns = [
        { title: 'Key', field: 'key', editable: "never"},
        { title: 'Value Used', field: 'value_used', editable:"never"},
        { title: 'Current Value', field: 'value' },
    ]
    const [data, setData] = useState([]);
    const [expanded, setExpanded] = React.useState(false);

    function reloadSetter(service_id, sr) {
        reload(service_id).then( results => {
            let d = []
            for (const [key, property] of Object.entries(sr["Properties"])) {
                let obj = {key: key, value_used: property["Value"], service_id: key}
                    results.forEach(res => {
                        if (key === res["key"]){
                            obj.value = res["value"]
                        }
                    })
                d.push(obj)
            }
            setData(d)
        })
    }

    const handleChange = (panel, service_id, service) => (event, isExpanded) => {
        setData([])
        if (!isExpanded){
            setExpanded(false)
        } else {
            setExpanded(panel);
            reloadSetter(service_id, service)
        }
    };
    const dt=props.dt
    const teamID = props.teamID
    return (
            <Box height="100%" width="100%" align="left" >
                {Object.keys(dt["Teams"][teamID]["Hosts"]).map((host) => {
                    let currentHost = dt["Teams"][teamID]["Hosts"][host]
                    return Object.keys(currentHost["Services"]).map((service_id) => {
                        let sr = currentHost["Services"][service_id]
                        let keyName
                        if (sr["DisplayName"]){
                            keyName = sr["DisplayName"]
                        } else {
                            if (currentHost["HostGroup"]){
                                keyName =currentHost["HostGroup"]["Name"] + "-" + sr["Name"]
                            } else{
                                keyName = sr["Name"]
                            }
                        }
                        return (
                            <Accordion expanded={expanded === keyName} onChange={handleChange(keyName, service_id, sr)} className={sr["Passed"] ? classes.customAccordionSuccessHeader: classes.customAccordionErrorHeader}>
                                <AccordionSummary
                                    expandIcon={<ExpandMoreIcon />}
                                    aria-controls={`${keyName}bh-content`}
                                    id={`${keyName}bh-header`}
                                >
                                    {sr["Passed"] ? <CheckCircleOutlineIcon className={classes.iconSuccess}  />  : <ErrorIcon className={classes.iconError}/>}
                                    <Typography className=  {classes.heading}>{keyName}</Typography>
                                    <Typography className={classes.secondaryHeading}>{}</Typography>

                                </AccordionSummary>
                                <AccordionDetails>
                                    {expanded === keyName &&
                                    <Box width="100%" bgcolor="background.paper">
                                        <Grid container spacing={3}>
                                            <Grid item xs={6}>
                                                {sr["Log"] &&
                                                <Alert severity={sr["Passed"] ? "info" : "warning"}>
                                                    <AlertTitle>Response</AlertTitle>
                                                    {sr["Log"]}
                                                </Alert>
                                                }
                                                <br/>
                                                {sr["Err"] &&
                                                <Alert severity="error">
                                                    <AlertTitle>Error Details</AlertTitle>
                                                    {sr["Err"]}
                                                </Alert>
                                                }
                                            </Grid>
                                            <Grid item xs={6}>

                                                <MaterialTable
                                                    options={{pageSizeOptions: [5,10,20,50,100], pageSize:20, emptyRowsWhenPaging:false}}
                                                    title="Properties"
                                                    columns={columns}
                                                    data={data}
                                                    cellEditable={{
                                                        onCellEditApproved: (newValue, oldValue, rowData, columnDef) => {
                                                            return new Promise((resolve) => {
                                                                setTimeout(() => {
                                                                    resolve();
                                                                    setData((prevState) => {
                                                                        return prevState.map(property => {
                                                                                if (property["key"] === rowData["key"]) {
                                                                                    return {...rowData, value: newValue}
                                                                                }
                                                                                return {...property}
                                                                            })
                                                                    });
                                                                    properties.Update(service_id, rowData["key"],{'value': newValue}).then( async () =>{
                                                                        reloadSetter(service_id, sr)
                                                                    }, (error) => {
                                                                        props.errorSetter(error)
                                                                    })
                                                                }, 600);
                                                            })
                                                        }
                                                    }}
                                                />

                                            </Grid>
                                        </Grid>
                                    </Box>
                                    }
                                </AccordionDetails>
                            </Accordion>
                        );
                    })
                })}
            </Box>
    );
}
