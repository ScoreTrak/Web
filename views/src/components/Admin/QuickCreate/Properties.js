import React, {useEffect} from 'react';
import {forwardRef, useImperativeHandle} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from "@material-ui/core/TextField";
import ServiceService from "../../../services/service/serivces"
import PropertiesService from "../../../services/property/properties"
import Box from "@material-ui/core/Box";
import CircularProgress from "@material-ui/core/CircularProgress";
import {Typography} from "@material-ui/core";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import PropertiesDefaults from "./serviceDefaultProperties"
import Switch from "@material-ui/core/Switch";
import Button from "@material-ui/core/Button";

const PropertiesCreate = forwardRef((props, ref) => {
    const [dt, setData] = React.useState({loader: true, services:[]})
    const [rowsData, setRowData] = React.useState({});


    useEffect(() => {
        ServiceService.GetAll().then(respService => {
            let rowdt = {}
            let displayNames = new Set()

            respService.forEach(serv => {
                if (serv.display_name){
                    displayNames.add(serv.display_name)
                    if (!(serv.display_name in rowdt)){
                        rowdt[serv.display_name] = {enableProcessingProperty: false}
                        Object.keys(PropertiesDefaults[serv.name]).forEach(key => {
                            rowdt[serv.display_name][key] = {
                                ...PropertiesDefaults[serv.name][key],
                                value: 'defaultValue' in PropertiesDefaults[serv.name][key] ? PropertiesDefaults[serv.name][key].defaultValue : '',
                                status: 'defaultStatus' in PropertiesDefaults[serv.name][key] ? PropertiesDefaults[serv.name][key].defaultStatus : 'View',
                            }
                        })
                    }
                }
            })

            setData({loader: false, services: respService})
            setRowData({...rowdt})
        }, props.errorSetter)
    }, []);


    const processProperties = (displayName, enable) => {
        setRowData(prevState => {return {...prevState, [displayName]: {...prevState[displayName], enableProcessingProperty: enable}}})
    }

    function submit() {
            let properties = []
            Object.keys(rowsData).forEach(DisplayName =>{

                if (rowsData[DisplayName].enableProcessingProperty){

                    dt.services.forEach(service => {
                        if (service.display_name === DisplayName){
                            Object.keys(rowsData[DisplayName]).forEach(propertyKey => {
                                if (propertyKey !== "enableProcessingProperty"){
                                    properties.push({
                                        service_id: service.id,
                                        key: propertyKey,
                                        status: rowsData[DisplayName][propertyKey].status,
                                        value: rowsData[DisplayName][propertyKey].value
                                    })
                                }
                            })
                        }
                    })
                }
            })
            props.handleLoading()
            PropertiesService.Create(properties).then(() => {props.handleSuccess()}, props.errorSetter)
    }

    const columns = [
        { id: 'key', label: 'Key'},
        { id: 'value', label: 'Value'},
        { id: 'status', label: 'Status'},
    ];

    return (
        <React.Fragment>
        <div>
            {!dt.loader ?
                Object.keys(rowsData).map((table) => (
                    <React.Fragment>
                        <Typography>Properties for: {table}</Typography>
                        <Table stickyHeader aria-label="sticky table" style={{marginBottom: '4vh'}}>
                            <TableHead>
                                <TableRow>
                                    {columns.map((column) => (
                                        <TableCell key={column.id} >
                                            {column.label}
                                        </TableCell>
                                    ))}
                                    <TableCell>
                                        <Switch
                                            checked={rowsData[table].enableProcessingProperty}
                                            onChange={(event) => {processProperties(table, event.target.checked)}}
                                            inputProps={{ 'aria-label': 'secondary checkbox' }}
                                        />
                                    </TableCell>
                                </TableRow>
                            </TableHead>

                            <TableBody>
                                { rowsData[table].enableProcessingProperty &&
                                    Object.keys(rowsData[table]).filter(property => {
                                        if (property === 'enableProcessingProperty'){
                                            return false
                                        }
                                        return true
                                    }).map(property => {
                                        return <TableRow hover role="checkbox" tabIndex={-1} >
                                            {columns.map(column => (
                                                <TableCell>
                                                    {
                                                        column.id === 'key' && rowsData[table][property].name
                                                    }

                                                    {
                                                        column.id === 'value' && rowsData[table][property].type === 'field' &&
                                                        <TextField value={rowsData[table][property].value}
                                                                   onChange={(event => {
                                                                       const val = event.target.value
                                                                       setRowData(prevState => {
                                                                           return {...prevState, [table]: {  ...prevState[table], [property]: {...prevState[table][property], value: val}}}
                                                                       })
                                                                   })}
                                                        />
                                                    }
                                                    {   column.id === 'value' && rowsData[table][property].type === 'select' &&
                                                            <Select
                                                                value={rowsData[table][property].value}
                                                                onChange={(event => {
                                                                    const val = event.target.value
                                                                    setRowData(prevState => {
                                                                        return {...prevState, [table]: {  ...prevState[table], [property]: {...prevState[table][property], value: val}}}
                                                                    })
                                                                })}
                                                            >
                                                                {
                                                                    rowsData[table][property].options.map(stat => {
                                                                        return <MenuItem value={stat}>{stat}</MenuItem>
                                                                    })
                                                                }
                                                            </Select>
                                                    }
                                                    {
                                                        column.id === 'status' &&
                                                        <Select
                                                            value={rowsData[table][property].status}
                                                            onChange={(event => {
                                                                const val = event.target.value
                                                                setRowData(prevState => {
                                                                    return {...prevState, [table]: {  ...prevState[table], [property]: {...prevState[table][property], status: val}}}
                                                                })
                                                            })}
                                                        >
                                                            {
                                                                ["View", "Edit", "Hide"].map(stat => {
                                                                    return <MenuItem value={stat}>{stat}</MenuItem>
                                                                })
                                                            }
                                                        </Select>
                                                    }
                                                </TableCell>
                                            ))}

                                            <TableCell />
                                        </TableRow>
                                    })
                                }
                            </TableBody>
                        </Table>
                    </React.Fragment>
                ))
                :
                <Box height="100%" width="100%" m="auto">
                    <CircularProgress  />
                </Box>
            }
        </div>
        <div style={{display: 'flex',  justifyContent: 'flex-end'}}>
            <Button onClick={submit} variant="contained" style={{ marginRight: '8px', marginTop: '8px'}}>Submit</Button>
        </div>
    </React.Fragment>
    );
})

export default PropertiesCreate;