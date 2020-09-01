import React from 'react';
import {forwardRef, useImperativeHandle} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from "@material-ui/core/TextField";
import TeamService from "../../../services/team/teams";
import {parse_index} from "./helper";
import Button from "@material-ui/core/Button";


const TeamCreate = forwardRef((props, ref) => {
    const [rows, setRows] = React.useState([]);

    const columns = [
        { id: 'name', label: 'Team Name'},
        { id: 'index', label: 'Index', field: <TextField onChange={(event => {
        setRows(parse_index(event.target.value).map(idx => {
                return {
                    name: <TextField required label="Team Name" id={`name_${idx}`}/> ,
                    index: idx
                }
            }
        ))})} id="filled-helperText" label="Index" helperText="This field is used to create host addresses. Ex: 1,2,4-15"/> },
    ];

    function submit() {
            const teams = rows.map(row =>{
                return {index: row["index"], name: document.getElementById(`name_${row["index"]}`).value}
            })
            props.handleLoading()
            TeamService.Create(teams).then(() => {props.handleSuccess()}, props.errorSetter)
        }


    return (
        <React.Fragment>
            <div>
                <Table stickyHeader aria-label="sticky table">
                    <TableHead>
                        <TableRow>
                            {columns.map((column) => (
                                <TableCell key={column.id} >
                                    {column.label}
                                </TableCell>
                            ))}
                        </TableRow>
                    </TableHead>
                    <TableHead>
                        <TableRow>
                            {columns.map((column) => (
                                <TableCell style={{minWidth: "300px"}}
                                    key={column.id}
                                >
                                    {column.field}
                                </TableCell>
                            ))}
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {rows.map((row) => {
                            return (
                                <TableRow hover role="checkbox" tabIndex={-1} key={row.code}>
                                    {columns.map((column) => {
                                        const value = row[column.id];
                                        return (
                                            <TableCell key={column.id} align={column.align}>
                                                {column.format && typeof value === 'number' ? column.format(value) : value}
                                            </TableCell>
                                        );
                                    })}
                                </TableRow>
                            );
                        })}
                    </TableBody>
                </Table>
            </div>
            <div style={{display: 'flex',  justifyContent: 'flex-end'}}>
                 <Button onClick={submit} variant="contained" style={{ marginRight: '8px', marginTop: '8px'}}>Submit</Button>
            </div>
    </React.Fragment>
    );
})

export default TeamCreate;