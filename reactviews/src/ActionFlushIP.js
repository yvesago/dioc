import { Button, useNotify } from 'react-admin';
import { useNavigate } from 'react-router-dom';
import { Alert } from '@mui/material';
import request from 'superagent';
import { MyConfig } from './MyConfig';


const ActionFlushIPButton = () => {
    const navigate = useNavigate();
    const notify = useNotify();
    const handleClick = () => {
        const token = localStorage.getItem('gotoken');
        request
            .put(MyConfig.API_URL + '/admin/api/v1/actionfluship') 
            .set('Authorization', `Bearer ${token}`)
            .then( (response)  => { 
                notify(<Alert severity="success">Flush IPs done</Alert>);
                // console.log(response);
                navigate('/extracts'); 
            })
            .catch((e) => {
                // console.error(e);
                notify(<Alert severity="error">Error in action Flush IPs</Alert>);
            });

    };

    return <Button label="Flush IPs" onClick={handleClick} />;
};


export default ActionFlushIPButton;
