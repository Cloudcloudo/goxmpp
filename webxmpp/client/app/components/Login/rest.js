import axios from 'axios';

const API = 'http://localhost:8080'

export function login(userData) {
  return axios.post(`${API}/login`, userData)
    .then(function (response) {
      return response.data;
    })
    .catch(function (error) {
        return false;
    });
}


// ASPMX.L.GOOGLE.COM/ALT1.ASPMX.L.GOOGLE.COM/ALT2.ASPMX.L.GOOGLE.COM/ALT3.ASPMX.L.GOOGLE.COM
