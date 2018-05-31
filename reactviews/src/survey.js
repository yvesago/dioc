import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, DisabledInput, required, EditButton, DateField, 
    RichTextField, SelectInput, Filter, Responsive, SimpleList
} from 'react-admin';
import RichTextInput from 'ra-input-rich-text';
import { withStyles } from '@material-ui/core/styles';
import { roles } from './MyConfig';

const levels = [
    { name: 'Critic', id: 'critic' },
    { name: 'Warn', id: 'warn' },
];

const colored = WrappedComponent => props => props.record.level === 'critic' ?
    <span style={{ color: 'red' }}><WrappedComponent {...props} /></span> :
    props.record.level === 'warn' ? 
        <span style={{ color: 'orange' }}><WrappedComponent {...props} /></span> :
        <WrappedComponent {...props} />;

const ColoredTextField = colored(TextField);


const SurveyFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} allowEmpty alwaysOn />
    </Filter>
);

const styles = {
    field: {
        width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const SurveyList = withStyles(styles)(({ classes, ...props }) => (
    <List filters={<SurveyFilter />} perPage={30} {...props}>
        <Responsive
            small={
                <SimpleList
                    primaryText={record => record.search}
                    secondaryText={record => `Role: ${record.role}, Level: ${record.level}`}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            }
            medium={
                <Datagrid>
                    <TextField source="search" />
                    <TextField source="role" />
                    <ColoredTextField source="level" />
                    <RichTextField source="comment" className={classes.field} stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            }
        />
    </List>
));


export const SurveyCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required()} />
            <SelectInput source="level" choices={levels} allowEmpty />
            <SelectInput source="role" choices={roles} allowEmpty />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);


export const SurveyEdit = (props) => (
    <Edit  {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required()} />
            <SelectInput source="role" choices={roles} allowEmpty />
            <SelectInput source="level" choices={levels} allowEmpty />
            <RichTextInput source="comment" />
            <DateField label="Created" source="created" showTime />
            <DateField label="Updated" source="updated" showTime />
            <DisabledInput source="crcs" validate={required()} />
        </SimpleForm>
    </Edit>
);
