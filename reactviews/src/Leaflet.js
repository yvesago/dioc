import React, { useEffect, useRef, useState } from 'react';
import L, { MarkerCluster } from 'leaflet';
import { MapContainer, useMapEvents, Marker, Popup, TileLayer } from 'react-leaflet';

import MarkerClusterGroup from 'react-leaflet-cluster';

import request from 'superagent';
import { MyConfig } from './MyConfig';

import { CircularProgress } from '@mui/material';
import { Box } from '@mui/material';

import 'leaflet/dist/leaflet.css';

function getLeafletPopup(name, ipname,loc) {
    var ipurl = MyConfig.BASE_PATH + '#/ips/' + name;
    return (
        <div>
            <b>{loc}:</b><br /> {ipname} »» <a href={ipurl}>View IP details</a>
        </div>
    );
}

const customIcon = new L.Icon({
    iconUrl: require('./location.svg').default,
    iconSize: new L.Point(40, 47)
});

const MapEventsHandler = ({ handleMapCursor }) => {
    useMapEvents({
        mousemove: (e) => handleMapCursor(e),
    });
    return null;
};


const MyLeaflet = (props) =>
{
    const latitude = props.lat || 46.76548;
    const longitude = props.lng || 12.76547;
    // Map state:
    const [map, setMap] = useState(null);
    const [markers, setMarkers] = useState(null);
    //console.log(props);

    const [coords, setCoords] = useState({});


    const handleMapCursor = (e) => {
        if ( props.point !== undefined ) {
            e.latlng && setCoords({ lat: e.latlng.lat, lng: e.latlng.lng});
        }
    };

    useEffect(() => {

        //console.log('useEffect');
        const a = 1;

        if ( props.point === undefined ) {
            const token = localStorage.getItem('gotoken');
            var u = MyConfig.API_URL + '/admin/api/v1/geojson';
            request
                .get(u)
                .set('Authorization', `Bearer ${token}`)
                .end(function(error, response){
                    if (error) return console.error(' --' + error);
                    var ms = [];
                    response.body.features.forEach(function(p) {
                        var m = {position:[p.geometry.coordinates[1], p.geometry.coordinates[0]],
                            popup:null,options:{}};
                        m.popup = getLeafletPopup( p.properties.Title || '',  p.properties.IP || '', p.properties.Loc || '');
                        m.options = { IP: p.properties.IP };
                        ms.push(m);
                    });
                    //console.log(ms);
                    setMarkers(ms);
                }.bind());
        } else {
            var ms = [];
            var m = {position:[props.lat,props.lng],popup:'point '+props.point, options: {IP: props.name}};
            m.popup = props.name; //getLeafletPopup('point '+props.point,props.name);
            //console.log(m);
            ms.push(m);
            setMarkers(ms);
        }

        if (!map) return;
    }, [map]);

    const { lat, lng } = coords;

    return (
        <>
            <MapContainer whenCreated={setMap} center={[latitude, longitude]} zoom={props.zoom} scrollWheelZoom={true} style={{height: '70vh' }}>
                <TileLayer
                    attribution='&copy; <a href="https://www.osm.org/copyright">OpenStreetMap</a> contributors'
                    url='https://{s}.tile.osm.org/{z}/{x}/{y}.png'
                />
                <MarkerClusterGroup chunkedLoading>
                    { markers && markers.map((o,index) =>  (
                        <Marker
                            key={index}
                            position={[o.position[0],o.position[1]]}
                            title={o.options.IP}
                            icon={customIcon}
                        >
                            <Popup>{o.popup}</Popup>
                        </Marker>
                    ))}
                />
                </MarkerClusterGroup>
                <MapEventsHandler handleMapCursor={handleMapCursor} />
                {/* Additional map layers or components can be added here */}
            </MapContainer>
            {
                <div>
                    <b>lat.</b>: {lat && lat.toFixed(4)} <b>lon.</b>: {lng && lng.toFixed(4)}
                </div>
            }
        </>
    );
};

export default MyLeaflet;
