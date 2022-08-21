/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
 let chai = require('chai');
 var expect = chai.expect
 let chaiHttp = require('chai-http');
 
 let should = chai.should();
 chai.use(chaiHttp);
 chai.use(require('chai-json'));
 console.log("running tests against host", process.env.HOST)
 describe('Users', function () {
   describe('/GET users', () => {
     it('it should GET all the users', (done) => {
       chai.request(process.env.HOST)
         .get('/users')
         .end((err, res) => {
           res.should.have.status(200);
           res.header["content-type"].should.be.eql('application/json')
           res.body.should.be.a('array');
           res.body.length.should.be.gt(0);
           done();
         });
     });
   });
 });