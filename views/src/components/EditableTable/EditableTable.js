import React, {useEffect} from "react";
import Paper from "@material-ui/core/Paper";
import Box from '@material-ui/core/Box';
import CircularProgress from "@material-ui/core/CircularProgress";
import MaterialTable from "material-table";


export default function Setup(props) {
    const setTitle = props.setTitle
    const classesPaper = props.classesPaper
    setTitle(props.title)

    const [state, setState] = React.useState({
        columns: props.columns,
        loader:true
    });
    async function reload() {
        if (props.isDependant){
            for (let owner = 0; owner < props.owningService.length; owner ++){
                let lookup = {}
                const owningObject = await props.owningService[owner].GetAll()
                for (let i = 0; i < owningObject.length; i++){
                    lookup[owningObject[i]["id"]] = `${owningObject[i][props.owningFieldLookup[owner]]}`
                    if (props.owningFieldLookup[owner] !== "id"){
                        lookup[owningObject[i]["id"]] += `(ID:${owningObject[i]["id"]})`
                    }
                }
                setState(prevState => {
                    let columns = prevState.columns
                    for (let i = 0; i < columns.length; i++){
                        if (columns[i]['field'] === props.fieldForLookup[owner]){
                            columns[i]['lookup'] = lookup
                        }
                    }
                    return{...prevState, columns: columns,
                    }})
                }
        }

        let objects = await props.service.GetAll()
        return {objects}


    }
    function reloadSetter() {
        reload().then(newState => { setState(prevState => {return{...prevState, data: newState.objects, loader:false}})}, props.errorSetter)
    }
    useEffect(() => {
        reloadSetter()
    }, []);

    return (
        <Paper className={classesPaper} style={{minHeight: "85vh"}}>
            <link
                rel="stylesheet"
                href="https://fonts.googleapis.com/icon?family=Material+Icons"
            />
            {!state.loader ?
                <Box height="100%" width="100%" align="left" >
                    <MaterialTable
                        title={props.title}
                        columns={state.columns}
                        data={state.data}
                        actions={(() => {
                            let arr = []
                            for (let i = 0; i < props.additionalActions.length; i++){
                                arr.push({...props.additionalActions[i], onClick: ((event, rowData) => { return props.additionalActions[i].onFuncClick(event, rowData).then( async () =>{
                                        reloadSetter()
                                    }, (error) => {
                                        reloadSetter()
                                        props.errorSetter(error)
                                    })})  })
                            }
                            return arr
                        })()}
                        options={{pageSizeOptions: [5,10,20,50,100, 500, 1000], pageSize:20, emptyRowsWhenPaging:false}}
                        editable={props.editable &&{
                            onRowAdd: (newData) =>
                                new Promise((resolve) => {
                                    setTimeout(() => {
                                        resolve();
                                        setState((prevState) => {
                                            const data = [...prevState.data];
                                            data.push(newData);
                                            return { ...prevState, data };
                                        });

                                        props.service.Create(props.disallowBulkAdd ? newData : [newData]).then( async () =>{
                                            reloadSetter()
                                        }, (error) => {
                                            reloadSetter()
                                            props.errorSetter(error)
                                        })
                                    }, 600);
                                }),
                            onRowUpdate: (newData, oldData) =>
                                new Promise((resolve) => {
                                    setTimeout(() => {
                                        resolve();
                                        if (oldData) {
                                            setState((prevState) => {
                                                const data = [...prevState.data];
                                                data[data.indexOf(oldData)] = newData;
                                                return { ...prevState, data };
                                            });

                                            let finalObj = {}

                                            for (let property in newData) {
                                                if (newData.hasOwnProperty(property)) {
                                                    if (oldData[property] !== newData[property]){
                                                        finalObj[property] = newData[property]
                                                    }
                                                }
                                            }
                                            props.service.Update(...props.idFields.map(id => {
                                                return oldData[id]
                                            }), finalObj).then( async () =>{
                                                reloadSetter()
                                            }, (error) => {
                                                reloadSetter()
                                                props.errorSetter(error)
                                            })
                                        }
                                    }, 600);
                                }),
                            onRowDelete: (oldData) =>
                                new Promise((resolve) => {
                                    setTimeout(() => {
                                        resolve();
                                        setState((prevState) => {
                                            const data = [...prevState.data];
                                            data.splice(data.indexOf(oldData), 1);
                                            return { ...prevState, data };
                                        });

                                        props.service.Delete(...props.idFields.map(id => {
                                            return oldData[id]
                                        })).then( async () =>{
                                            reloadSetter()
                                        }, (error) => {
                                            reloadSetter()
                                            props.errorSetter(error)
                                        })
                                    }, 600);
                                }),
                        }}
                    />
                </Box>
                :
                <Box height="100%" width="100%" m="auto">
                    <CircularProgress  />
                </Box>
            }

        </Paper>
    );
}