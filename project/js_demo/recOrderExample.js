const FabricCAServices = require('fabric-ca-client');
const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

// Load the network configuration
const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

// Create a new CA client for interacting with the CA.
const caInfo = ccp.certificateAuthorities['ca.org1.example.com'];
const caTLSCACerts = caInfo.tlsCACerts.pem;
const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

// Create a new file system based wallet for managing identities.
const walletPath = path.join(process.cwd(), 'wallet');

test();

async function test() {
    console.log("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ Start Test ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++");

    // Test enroll and register
    let userName = "testUser";
    await enrollAdmin();
    await registerUser(userName);

    // Initialize the ledger
    await invoke("InitLedger");

    // Test query
    let queryType = "user";
    let queryID = "user1";
    await queryAll(queryType);
    await queryByID(queryType, queryID);

    queryType = "recOrder"
    queryID = "order2"
    await queryAll(queryType);
    await queryByID(queryType, queryID);

    queryType = "receivable"
    queryID = "rec3"
    await queryAll(queryType);
    await queryByID(queryType, queryID);

    // Test sample routine
    let company = "user1";
    let firstSupplier = "user2";
    let secondSupplier = "user3";
    let financial = "user4"

    console.log("The company is creating a receivable order...");
    let newOrder = await invoke("CreateRecOrder", company, firstSupplier, 500000);
    console.log(newOrder);
    console.log('========================================================================')

    console.log("The first supplier is signing the receivable order...")
    let rec = await invoke("SignReceivable", firstSupplier, newOrder.order_no, 499999);
    console.log(rec);
    console.log('========================================================================')

    console.log("The company is accepting the receivable...")
    let rec = await invoke("AcceptReceivable", company, rec.receivable_no);
    console.log(rec);
    console.log('========================================================================')

    console.log("The first supplier is transferring the receivable...")
    rec = await invoke("TransferReceivable",  firstSupplier, secondSupplier, rec.receivable_no,)
    console.log(rec);
    console.log('========================================================================')

    console.log("The second supplier is applying for a discount from the financial...")
    rec = await invoke("ApplyDiscount", secondSupplier, financial, rec.receivable_no);
    console.log(rec);
    console.log('========================================================================')

    console.log("The financial is confirming discount application from the second supplier...")
    rec = await invoke("DiscountConfirm", financial, rec.order_no);
    console.log(rec);
    console.log('========================================================================')

    console.log("The company is paying the account of receivable...")
    rec = await invoke("Redeemed", company, rec.order_no);
    console.log(rec);

    console.log("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ End Test ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++");
    process.exit(1);
}

async function enrollAdmin() {
    try {
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get('admin');
        if (identity) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            console.log('========================================================================');
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: 'Org1MSP',
            type: 'X.509',
        };
        await wallet.put('admin', x509Identity);
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
        console.log('========================================================================');
    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
        console.log('========================================================================');
        process.exit(1);
    }
}

async function registerUser(userName) {
    try {
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userIdentity = await wallet.get(userName);
        if (userIdentity) {
            console.log('An identity for the user "' + userName + '" already exists in the wallet');
            console.log('========================================================================');
            return;
        }

        // Check to see if we've already enrolled the admin user.
        const adminIdentity = await wallet.get('admin');
        if (!adminIdentity) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            console.log('========================================================================');
            return;
        }

        // build a user object for authenticating with the CA
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, 'admin');

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register({
            affiliation: 'org1.department1',
            enrollmentID: userName,
            role: 'client'
        }, adminUser);
        const enrollment = await ca.enroll({
            enrollmentID: userName,
            enrollmentSecret: secret
        });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: 'Org1MSP',
            type: 'X.509',
        };
        await wallet.put(userName, x509Identity);
        console.log('Successfully registered and enrolled admin user "' + userName + '" and imported it into the wallet');
        console.log('========================================================================');
    } catch (error) {
        console.error(`Failed to register user: ${error}`);
        console.log('========================================================================');
        process.exit(1);
    }
}

async function queryAll(queryType) {
    try {
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("admin");
        if (!identity) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('========================================================================');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: "admin", discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('rec_order');

        let result;
        if (queryType === "user") {
            result = await contract.evaluateTransaction('QueryAllUsers');
        } else if (queryType === "recOrder") {
            result = await contract.evaluateTransaction('QueryAllRecOrders');
        } else if (queryType === "receivable"){
            result = await contract.evaluateTransaction('QueryAllReceivables');
        } else {
            console.error("No such type for querying");
            process.exit(1);
        }
        let res = JSON.parse(result.toString());
        console.log('Transaction has been evaluated, result is:');
        console.log(res)
        console.log('========================================================================');
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

async function queryByID(queryType, ID) {
    try {
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("admin");
        if (!identity) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('========================================================================');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: "admin", discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('rec_order');

        let result;
        if (queryType === "user") {
            result = await contract.evaluateTransaction("QueryUser", ID);
        } else if (queryType === "recOrder") {
            result = await contract.evaluateTransaction('QueryRecOrder', ID);
        } else if (queryType === "receivable") {
            result = await contract.evaluateTransaction('QueryReceivable', ID);
        } else {
            console.error("No such type for querying");
            process.exit(1);
        }
        let res = JSON.parse(result.toString());
        console.log('Transaction has been evaluated, result is:');
        console.log(res)
        console.log('========================================================================');
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

async function invoke(methodName, ...args) {
    let res;
    try {
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get('admin');
        if (!identity) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('========================================================================');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('rec_order');

        let result = await contract.submitTransaction(methodName, ...args);
        if (methodName === "InitLedger") {
            return null;
        }
        console.log('Transaction has been submitted');
        res = JSON.parse(result.toString());

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        console.log('========================================================================');
        process.exit(1);
    }
    return res;
}