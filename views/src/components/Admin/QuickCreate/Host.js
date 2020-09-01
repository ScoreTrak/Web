import React, {useEffect} from 'react';
import {forwardRef, useImperativeHandle} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from "@material-ui/core/TextField";
import TeamService from "../../../services/team/teams";
import HostGroupService from "../../../services/host_group/host_groups"
import HostService from "../../../services/host/hosts"
import Box from "@material-ui/core/Box";
import CircularProgress from "@material-ui/core/CircularProgress";
import Button from "@material-ui/core/Button";


const HostCreate = forwardRef((props, ref) => {

    const [dt, setData] = React.useState({loader: true, teams:[], hostGroups:[]})

    useEffect(() => {
        TeamService.GetAll().then(respTeam => { HostGroupService.GetAll().then(respHostGroup => { setData({teams: respTeam.sort((a, b) => (a.index > b.index) ? 1 : -1), hostGroups: respHostGroup, loader: false}) }, props.errorSetter)}, props.errorSetter)
    }, []);

    const [rowsData, setRowData] = React.useState({});

    function modifyTeamRows(host_group_id, templateValue){
        let nextRowData = {}
        if (templateValue.includes('X')){
            for (let i = 0; i < dt.teams.length; i++){
                if (dt.teams[i].index) {
                    nextRowData[`${dt.teams[i].id}_${host_group_id}`] = templateValue.replace("X", dt.teams[i].index)
                }
            }

            setRowData(prevState => {return{...prevState, ...nextRowData}})
        }
    }


    function submit() {
            let hosts = []
            Object.keys(rowsData).forEach(team_hostGrp_id => {
                let team_id, host_group_id
               [team_id, host_group_id] = team_hostGrp_id.split("_")
                if (rowsData[team_hostGrp_id]){
                    hosts.push({
                        team_id: team_id,
                        host_group_id: host_group_id,
                        address: rowsData[team_hostGrp_id],
                        enabled: true
                    })
                }
            })
            props.handleLoading()
            HostService.Create(hosts).then(() => {props.handleSuccess()}, props.errorSetter)
    }




    return (
        <React.Fragment>
        <div>
            {!dt.loader ?
            <Table stickyHeader aria-label="sticky table">
                <TableHead>
                    <TableRow>
                        <TableCell />
                        {dt.hostGroups.map((column) => (
                            <TableCell>
                                {column.name}
                            </TableCell>
                        ))}
                    </TableRow>
                </TableHead>
                <TableHead>
                    <TableRow>
                        <TableCell />
                        {dt.hostGroups.map((column) => (
                            <TableCell>
                                <TextField label="Template" id={`id_${column.id}`} helperText="Ex. 10.1.X.1" onChange={event => {modifyTeamRows(column.id, event.target.value)}}/>
                            </TableCell>
                        ))}
                    </TableRow>
                </TableHead>
                <TableBody>
                    {dt.teams.map((row) => {
                        if (row.index){
                            return (
                                <TableRow hover role="checkbox" tabIndex={-1} key={row.code}>
                                    <TableCell key={row["name"]}>
                                        {row["name"]}
                                    </TableCell>

                                    {dt.hostGroups.map((column) => {
                                        return (
                                            <TableCell>
                                                <TextField id={`${row.id}_${column.id}`} value={rowsData[`${row.id}_${column.id}`]} onChange={(event => {
                                                    const val = event.target.value
                                                    setRowData(prevState => {
                                                        return {...prevState, [`${row.id}_${column.id}`]: val}
                                                    })
                                                })}
                                                />
                                            </TableCell>
                                        );
                                    })}
                                </TableRow>
                            );
                        }
                    })}
                </TableBody>
            </Table>
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

export default HostCreate;