import React, {useEffect} from 'react';
import {forwardRef, useImperativeHandle} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from "@material-ui/core/TextField";
import HostService from "../../../services/host/hosts";
import HostGroupService from "../../../services/host_group/host_groups"
import ServiceGroupService from "../../../services/service_group/service_groups"
import ServiceService from "../../../services/service/serivces"
import Box from "@material-ui/core/Box";
import CircularProgress from "@material-ui/core/CircularProgress";
import {Typography} from "@material-ui/core";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import PropertiesDefaults from "./serviceDefaultProperties"
import Button from "@material-ui/core/Button";

const ServiceCreate = forwardRef((props, ref) => {
    const [dt, setData] = React.useState({loader: true, hosts:[], hostGroups:[], serviceGroups: []})
    const [counter, setCounter] = React.useState({})
    const [rowsData, setRowData] = React.useState({});

    const columns = [
        { id: 'name', label: 'Name of the Check(Ex: PING, SSH)'},
        { id: 'display_name', label: 'Display Name'},
        { id: 'points', label: 'Points'},
        { id: 'points_boost', label: 'Points Boost'},
        { id: 'round_units', label: 'Round Units'},
        { id: 'round_units_delay', label: 'Round Delay'},
        { id: 'service_group_id', label: 'Service Group'}
    ];

    const defaultVals = {
        name: '',
        display_name: '',
        points: 1,
        'points_boost': 0,
        'round_units': 1,
        'round_units_delay': 0,
        'service_group_id' : ''
    }

    useEffect(() => {
        ServiceGroupService.GetAll().then(respServiceGrp => { HostGroupService.GetAll().then(respHostGroup => { HostService.GetAll().then(respHost => {
            let counter = {}
            let rowdt = {}
            respHostGroup.forEach(hstGrp => {
                counter[hstGrp.id] = 1
                rowdt[hstGrp.id] = {1: defaultVals}
            })
            setCounter({...counter})
            setData({loader: false, hosts: respHost, hostGroups: respHostGroup, serviceGroups: respServiceGrp})
            setRowData({...rowdt})


        }, props.errorSetter)}, props.errorSetter)}, props.errorSetter)
    }, []);

    const setNumberOfServices = (hostGroupID, value) => {
        if (value === ""){
            return
        }
        if (value >= 0){
            setRowData(prevState => {
                let newRowData = {}
                for (let i = 1; i <= value; i++){
                    if (i in prevState[hostGroupID]){
                        newRowData[i] = prevState[hostGroupID][i]
                    } else {
                        newRowData[i] = defaultVals
                    }
                }
                return {...prevState, [hostGroupID]: newRowData}
            })
            setCounter(prevState => {return {...prevState, [hostGroupID]: value}})
        }
    }


    function submit() {
            let services = []
            Object.keys(rowsData).forEach(hostGroupID =>{
                Object.keys(rowsData[hostGroupID]).forEach((idx) => {
                    dt.hosts.forEach(host => {
                        if (hostGroupID === host.host_group_id){
                            services.push({...rowsData[hostGroupID][idx], host_id: host.id})
                        }
                    })
                })
            })
            props.handleLoading()
            ServiceService.Create(services).then(() => {props.handleSuccess()}, props.errorSetter)
        }

    return (
        <React.Fragment>
        <div>
            {!dt.loader ?
                dt.hostGroups.map((table) => (
                        <React.Fragment>
                            <Typography>Services for: {table.name}</Typography>
                            <Table stickyHeader aria-label="sticky table" style={{marginBottom: '4vh'}}>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>
                                            <TextField label="#" helperText="Number of services" type="number" value={counter[table.id]} onChange={event => {setNumberOfServices(table.id, event.target.value)}}/>
                                        </TableCell>
                                        {columns.map((column) => (
                                            <TableCell key={column.id} >
                                                {column.label}
                                            </TableCell>
                                        ))}
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {
                                        Array.apply(null, { length: counter[table.id] }).map((e,j) => {
                                            let i = j+1
                                            return <TableRow hover role="checkbox" tabIndex={-1} key={`${table.id}_${i}`}>
                                                    <TableCell key={`${table.id}_${i}`}>
                                                        {i}
                                                    </TableCell>

                                                    {columns.map(column => (
                                                         <TableCell>
                                                             {
                                                                 column.id === 'name' &&
                                                                 <Select
                                                                     id={`${table.id}_${i}_${column.id}`}
                                                                     value={rowsData[table.id] && rowsData[table.id][i] && rowsData[table.id][i][column.id]}
                                                                     onChange={(event => {
                                                                         const val = event.target.value
                                                                         setRowData(prevState => {
                                                                             return {...prevState, [table.id]: {...prevState[table.id], [i]: {...prevState[table.id][i], [column.id]: val}}}
                                                                         })})}
                                                                 >
                                                                     {

                                                                         Object.keys(PropertiesDefaults).map(serviceName => {
                                                                             return <MenuItem value={serviceName}>{serviceName}</MenuItem>
                                                                         })

                                                                     }
                                                                 </Select>
                                                             }

                                                             {
                                                                 column.id === 'service_group_id' &&
                                                                 <Select
                                                                     id={`${table.id}_${i}_${column.id}`}
                                                                     value={rowsData[table.id] && rowsData[table.id][i] && rowsData[table.id][i][column.id]}
                                                                     onChange={(event => {
                                                                         const val = event.target.value
                                                                         setRowData(prevState => {
                                                                             return {...prevState, [table.id]: {...prevState[table.id], [i]: {...prevState[table.id][i], [column.id]: val}}}
                                                                         })})}
                                                                 >
                                                                     {
                                                                         dt.serviceGroups.map(servGroup => {
                                                                             return <MenuItem value={servGroup.id}>{servGroup.name}</MenuItem>
                                                                         })
                                                                     }
                                                                 </Select>
                                                             }

                                                             {
                                                                 (column.id !== 'service_group_id' && column.id !== 'name') &&
                                                                     <TextField id={`${table.id}_${i}_${column.id}`} type={(column.id==='points' || column.id==='points_boost' || column.id==='points' || column.id === 'round_units' || column.id === 'round_units_delay') && 'number'} value={rowsData[table.id] && rowsData[table.id][i] && rowsData[table.id][i][column.id]}
                                                                         onChange={(event => {
                                                                             let val = event.target.value
                                                                             if ((column.id==='points' || column.id==='points_boost' || column.id==='points' || column.id === 'round_units' || column.id === 'round_units_delay')){
                                                                                 val = parseInt(val)
                                                                             }
                                                                             setRowData(prevState => {
                                                                                 return {...prevState, [table.id]: {...prevState[table.id], [i]: {...prevState[table.id][i], [column.id]: val}}}
                                                                         })})}
                                                                                //Todo: Validate on positive numbers
                                                                     />
                                                             }
                                                         </TableCell>
                                                    ))}
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

export default ServiceCreate;