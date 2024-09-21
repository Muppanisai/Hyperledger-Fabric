'use strict';

const { Contract } = require('fabric-contract-api');

class AssetContract extends Contract {

    // CreateAsset adds a new asset to the world state
    async CreateAsset(ctx, dealerID, msisdn, mpin, balance, status) {
        const asset = {
            dealerID,
            msisdn,
            mpin,
            balance: parseFloat(balance),
            status,
            transAmount: 0,
            transType: '',
            remarks: ''
        };

        await ctx.stub.putState(msisdn, Buffer.from(JSON.stringify(asset)));
        return `Asset ${msisdn} created successfully`;
    }

    // UpdateBalance modifies the balance of an existing asset
    async UpdateBalance(ctx, msisdn, transAmount, transType, remarks) {
        const assetAsBytes = await ctx.stub.getState(msisdn); // Get the asset from the ledger
        if (!assetAsBytes || assetAsBytes.length === 0) {
            throw new Error(`Asset ${msisdn} does not exist`);
        }

        const asset = JSON.parse(assetAsBytes.toString());

        if (transType === 'debit') {
            asset.balance -= parseFloat(transAmount);
        } else if (transType === 'credit') {
            asset.balance += parseFloat(transAmount);
        } else {
            throw new Error('Transaction type must be either debit or credit');
        }

        asset.transAmount = parseFloat(transAmount);
        asset.transType = transType;
        asset.remarks = remarks;

        await ctx.stub.putState(msisdn, Buffer.from(JSON.stringify(asset)));
        return `Balance for asset ${msisdn} updated successfully`;
    }

    // QueryAsset returns the asset details based on the MSISDN
    async QueryAsset(ctx, msisdn) {
        const assetAsBytes = await ctx.stub.getState(msisdn);
        if (!assetAsBytes || assetAsBytes.length === 0) {
            throw new Error(`Asset ${msisdn} does not exist`);
        }

        return assetAsBytes.toString();
    }

    // GetAssetHistory returns the transaction history for a given asset
    async GetAssetHistory(ctx, msisdn) {
        const iterator = await ctx.stub.getHistoryForKey(msisdn);
        const history = [];

        while (true) {
            const res = await iterator.next();

            if (res.value && res.value.value.toString()) {
                const jsonRes = {};
                jsonRes.txId = res.value.tx_id;
                jsonRes.timestamp = res.value.timestamp;
                jsonRes.isDelete = res.value.is_delete.toString();
                jsonRes.value = JSON.parse(res.value.value.toString('utf8'));
                history.push(jsonRes);
            }

            if (res.done) {
                await iterator.close();
                return JSON.stringify(history);
            }
        }
    }
}

module.exports = AssetContract;
