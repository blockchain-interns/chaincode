//ROHAN NODEJS CHAINCODE
//CHAINCODE

'use strict';

//IMPORTS
const shim = require('fabric-shim');
const util = require('util');

let Chaincode = class {

//The Init method is called when this chaincode is instantiated by the blockchain network
//Ledger instantiation is being handled by a separate function - see initLedger()
async Init(stub) {
	console.info('=========== Instantiated Rohan\'s chaincode ===========');
	return shim.success();
}

  //The Invoke method is called when an application requests to run this
  //chaincode. The calling application program must also specify the function 
  //to be called (fcn), with the correct arguments.
  async Invoke(stub) {
    let ret = stub.getFunctionAndParameters();
    console.info(ret);

    //Set 'method' to whatever function is requested by the app
    let method = this[ret.fcn]; 
    if (!method) {
      console.error('The function '+ ret.fcn + ' was not found.');
      throw new Error('Received unknown function - ' + ret.fcn + ' - invocation');
    }
    //Ensure the correct arguments were passed by the application 
    try {
      let payload = await method(stub, ret.params);
      return shim.success(payload);
    } catch (err) {
      console.log(err);
      return shim.error(err);
    }
  }

  //initLedger() - Instantiates the blockchain ledger with some default objects
  //In this case, adding student objects with attributes: {name, age, grade, year}
  async initLedger(stub, args) {
    console.info('============= START : Initialize Ledger ===========');
     //Create array of student objects
     let students = [];
     students.push({
          name: 'Emily',
          age: '16',
          grade: '76',
          year: '2'
      });
      students.push({
          name: 'John',
          age: '17',
          grade: '80',
          year: '3'
      });
      students.push({
          name: 'Ashley',
          age: '15',
          grade: '89',
          year: '1'
      });
      students.push({
          name: 'Rohan',
          age: '18',
          grade: '100',
          year: '4'
      });

    //Add students to ledger
    for (let i = 0; i < students.length; i++) {
      students[i].docType = 'student';
        await stub.putState('STUDENT' + i, Buffer.from(JSON.stringify(students[i])));
        console.info(`\nStudent - ${students[i].name} - was added to the ledger.\n`);
    }
    console.info('============= END : Initialize Ledger ===========');
  }

  //addStudent() - Adds a student object to the ledger
  async addStudent(stub, args) {
    console.info('============= START : Add Student ===========');
      //Check that 5 arguments were passed by app (one for docType, 4 for student object)
      if (args.length != 5) {
      throw new Error('Incorrect number of arguments. Expecting 5.');
    }

    //Create student object
    var student = {
      docType: 'student',
      name: args[1],
      age: args[2],
      grade: args[3],
      year: args[4]
    };

    //Add student to the ledger
    await stub.putState(args[0], Buffer.from(JSON.stringify(student)));
    console.info(`\nStudent - ${student.name} - was added to the ledger.\n`);
    console.info('============= END : Create Car ===========');
  }

    //queryStudent() - Queries the ledger for a particular student by ID
    async queryStudent(stub, args) {
        if (args.length != 1) {
            throw new Error('Invalid argument or too many arguments passed. Expecting a student ID. Example: STUDENT3');
        }
        let studentName = args[0];

        let studentAsBytes = await stub.getState(studentName); //Get the student from the chaincode state
        if (!studentAsBytes || studentAsBytes.toString().length <= 0) {
            throw new Error(studentName + ' does not exist: ');
        }
        console.log(studentAsBytes.toString());
        return studentAsBytes;
    }

  //queryAllStudents() - Queries the ledger for all information
  async queryAllStudents(stub, args) {

    //Sets a range to query based on keys that were specified in initLedger() by the stub.putState() function
    let startKey = 'STUDENT0';
    let endKey = 'STUDENT999';

    let iterator = await stub.getStateByRange(startKey, endKey);

    console.info('============= START: Querying All Students ===========\n');

    let allResults = [];
    while (true) {
      let res = await iterator.next();

      if (res.value && res.value.value.toString()) {
        let jsonRes = {};
        console.log(res.value.value.toString('utf8'));

        jsonRes.Key = res.value.key;
        try {
          jsonRes.Record = JSON.parse(res.value.value.toString('utf8'));
        } catch (err) {
          console.log(err);
          jsonRes.Record = res.value.value.toString('utf8');
        }
        allResults.push(jsonRes);
      }
      if (res.done) {
        console.log('End of data');
        await iterator.close();
        console.info(allResults);
        return Buffer.from(JSON.stringify(allResults));
      }
    }
      console.info('\n============= END: Querying All Students ===========');
  }

    //deleteStudent() - Removes a student object from the ledger
    async deleteStudent(stub, args) {
        console.info('============= START : Delete Student ===========');
        if (args.length != 1) {
            throw new Error('Invalid or incorrect number of arguments passed. Expecting Student ID. Example: STUDENT2');
    }

        let studentAsBytes = await stub.getState(args[0]);

        await stub.deleteState(studentAsBytes);
        console.info('============= END : Delete Student ===========');

  }
};

shim.start(new Chaincode());
