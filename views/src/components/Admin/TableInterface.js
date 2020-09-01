import EditableTable from "../EditableTable/EditableTable";
import React from "react";

export function Table(props, title, isDependant, columns, disallowBulkAdd, service, owningService, fieldForLookup, owningFieldLookup=["id"], editable=true, additionalActions=[], idFields=["id"]){
    return <EditableTable {...props} title={title} isDependant={isDependant} columns={columns} service={service} disallowBulkAdd={disallowBulkAdd} owningService={owningService}
                          fieldForLookup={fieldForLookup} owningFieldLookup={owningFieldLookup} editable={editable} additionalActions={additionalActions} idFields={idFields}
    />
}
