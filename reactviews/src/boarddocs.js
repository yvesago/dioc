import React from 'react';
import { Edit, SimpleForm  } from 'react-admin';
import RichTextInput from 'ra-input-rich-text';
import Divider from '@material-ui/core/Divider';

export const BoardEdit = (props) => (
    <Edit title="Docs" undoable={false} {...props}>
        <SimpleForm redirect="/">
            <RichTextInput
                source="docs"
                label="Main Doc"
            />
            <br />
            <br />
            <Divider inset />
            <RichTextInput
                source="docagents"
                label="Agents"
            />
            <br />
            <br />
            <Divider inset />
            <RichTextInput
                source="docsearchs"
                label="Searchs"
            />
        </SimpleForm>
    </Edit>
);
