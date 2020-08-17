import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import AuthService from "../../services/auth/auth";


const useStyles = makeStyles({
    root: {
        width: '100%',
    },
    container: {
        maxHeight: 440,
    },
    tableNavigator:{
        marginRight: "5vh",
        marginLeft: "5vh"
    }
});

export default function EditableTable(props) {
    const classes = useStyles();
    const [rowPage, setRowPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(10);
    const [dense, setDense] = React.useState(false);
    const [hideAddresses, setHideAddresses] = React.useState(false);

    const handleHideAddresses = (event) => {
        setHideAddresses(event.target.checked);
    };
    const handleRowChangePage = (event, newPage) => {
        setRowPage(newPage);
    };
    const handleChangeRowsPerPage = (event) => {
        setRowsPerPage(+event.target.value);
        setRowPage(0);
    };
    
    const [columnPage, setColumnPage] = React.useState(0);
    const [columnsPerPage, setColumnsPerPage] = React.useState(10);
    const handleColumnChangePage = (event, newPage) => {
        setColumnPage(newPage);
    };
    const handleChangeColumnsPerPage = (event) => {
        setColumnsPerPage(+event.target.value);
        setColumnPage(0);
    };

    const handleChangeDense = (event) => {
        setDense(event.target.checked);
    };

    const dt = props.dt
    let teamNames = []
    let data = {}
    let dataKeys = new Set();
    if ("Teams" in dt){
        for (let team in dt["Teams"]) {
            if (dt["Teams"].hasOwnProperty(team)) {
                for (let host in dt.Teams[team]["Hosts"]){
                    if (dt.Teams[team]["Hosts"].hasOwnProperty(host)) {
                        let serviceAggregator = {}
                        if (Object.keys(dt.Teams[team]["Hosts"][host]["Services"]).length !== 0){
                            for (let service in dt.Teams[team]["Hosts"][host]["Services"]) {
                                if (dt.Teams[team]["Hosts"][host]["Services"].hasOwnProperty(service)) {
                                    let sr = dt.Teams[team]["Hosts"][host]["Services"][service]
                                    let keyName = ""
                                    if (sr["DisplayName"]){
                                        keyName = sr["DisplayName"]
                                    } else {
                                        if (dt.Teams[team]["Hosts"][host]["HostGroup"]){
                                            keyName = dt.Teams[team]["Hosts"][host]["HostGroup"]["Name"] + "-" + sr["Name"]
                                        } else{
                                            keyName = sr["Name"]
                                        }
                                    }
                                    serviceAggregator[keyName] = sr
                                    dataKeys.add(keyName)
                                }
                            }
                            teamNames.push(dt["Teams"][team]["Name"])
                            data[dt["Teams"][team]["Name"]] = {service: serviceAggregator, Address:dt["Teams"][team]["Hosts"][host]["Address"]}
                        }
                    }
                }
            }
        }
    }
    let dataKeysArray = [...dataKeys]
    return (
        <Paper className={classes.root}>
            <TableContainer className={classes.container}>
                <div>
                <Table stickyHeader aria-label="sticky table" size={dense ? 'small' : 'medium'}>
                    <TableHead>
                        <TableRow>
                            <TableCell
                                key="name"
                            >
                                Team Name
                            </TableCell>

                            {dataKeysArray.slice(columnPage * columnsPerPage, columnPage * columnsPerPage + columnsPerPage).map((column) => (
                                <TableCell align="center"
                                    key={column}
                                >
                                    {column}
                                </TableCell>
                            ))}
                        </TableRow>
                    </TableHead>


                    <TableBody>
                        {teamNames.slice(rowPage * rowsPerPage, rowPage * rowsPerPage + rowsPerPage).map((name) => {
                            return (
                                <TableRow hover tabIndex={-1} key={name}>
                                    <TableCell key={name}>
                                        {name}
                                    </TableCell>
                                    {dataKeysArray.slice(columnPage * columnsPerPage, columnPage * columnsPerPage + columnsPerPage).map((column) => (
                                        <TableCell key={name+column} style={(() => {
                                            if (data[name]["service"][column]) {
                                                if (data[name]["service"][column]["Passed"]){
                                                    return {backgroundColor: "green"}
                                                }
                                                return {backgroundColor: "red", color: "white"}
                                            }
                                        })()} align="center"
                                        >
                                            {!hideAddresses && data[name]["service"][column] && (() => {
                                                let msg = ""
                                                if (data[name]["Address"]) {
                                                    msg += data[name]["Address"]
                                                    if (column in data[name]["service"] && "Properties" in data[name]["service"][column]
                                                        && "Port" in data[name]["service"][column]["Properties"]) {
                                                        msg += ":" + data[name]["service"][column]["Properties"]["Port"]
                                                    }
                                                }
                                                return msg
                                            })()}
                                        </TableCell>
                                    ))}
                                </TableRow>
                            );
                        })}
                        
                    </TableBody>
                </Table>
                </div>
            </TableContainer>
            <div style={{display:"flex", flexDirection: "row", justifyContent: "center", alignItems: "center"}}>
                <TablePagination className={classes.tableNavigator}
                    rowsPerPageOptions={[1, 5, 10, 25, 100]}
                    component="div"
                    count={teamNames.length}
                    rowsPerPage={rowsPerPage}
                    page={rowPage}
                    onChangePage={handleRowChangePage}
                    onChangeRowsPerPage={handleChangeRowsPerPage}
                />
                <FormControlLabel className={classes.tableNavigator}
                                  control={<Switch checked={dense} onChange={handleChangeDense} />}
                                  label="Dense padding"
                />
                <FormControlLabel className={classes.tableNavigator}
                                  control={<Switch checked={hideAddresses} onChange={handleHideAddresses} />}
                                  label={"Hide Addresses"}
                />
                <TablePagination
                    labelRowsPerPage="Columns per page"
                    rowsPerPageOptions={[1, 5, 10, 25, 100]}
                    component="div"
                    count={dataKeysArray.length}
                    rowsPerPage={columnsPerPage}
                    page={columnPage}
                    className={classes.tableNavigator}
                    onChangePage={handleColumnChangePage}
                    onChangeRowsPerPage={handleChangeColumnsPerPage}
                />
            </div>
        </Paper>
    );
}