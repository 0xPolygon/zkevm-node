package state

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
)

// ProcessBlobInnerProcessRequest is the request to process a blob
// you must use the builder to create the request
type ProcessBlobInnerProcessRequest struct {
	oldBlobStateRoot    common.Hash
	oldBlobAccInputHash common.Hash
	oldNumBlob          uint64
	oldStateRoot        common.Hash
	forkId              uint64
	lastL1InfoTreeIndex uint32
	lastL1InfoTreeRoot  common.Hash
	timestampLimit      uint64
	coinbase            common.Address
	zkGasLimit          uint64
	blobType            BlobType
}

type ProcessBlobInnerProcessRequestBuilder struct {
	forkId      uint64
	blobType    BlobType
	isFirstBlob *bool
	err         error
	data        ProcessBlobInnerProcessRequest
}

// NewProcessBlobInnerProcessRequestBuilder creates a builder for ProcessBlobInnerProcessRequest
// The basic required information is:
//   - forkid: because in the future some fields could be different depending on the fork
//   - blobType: the type of the blob. That change the fields that are required
//     The class implement a fluid interface so a example of usage:
//     processRequest, err := NewProcessBlobInnerProcessRequestBuilder(10, TypeCallData).
//     SetAsFirstBlob().Build()

func NewProcessBlobInnerProcessRequestBuilder(forkid uint64, blob *BlobInner,
	previousSequence *BlobSequence,
	currentSequence BlobSequence) *ProcessBlobInnerProcessRequestBuilder {
	res := &ProcessBlobInnerProcessRequestBuilder{
		forkId:   forkid,
		blobType: BlobType(blob.Type),
		err:      nil,
		data: ProcessBlobInnerProcessRequest{
			forkId:           forkid,
			blobType:         BlobType(blob.Type),
			oldBlobStateRoot: ZeroHash, // Is always zero!
		},
	}
	if previousSequence == nil {
		res.setAsFirstBlob()
	} else {
		res.setPreviousSequence(*previousSequence)
	}
	res.setBlob(blob)
	res.setCurrentSequence(currentSequence)
	return res
}

func (p *ProcessBlobInnerProcessRequestBuilder) setAsFirstBlob() {
	if p.isFirstBlob != nil && !*p.isFirstBlob {
		p.err = errors.New("the blob is not the first blob, you can't set as first blob")
	}
	tmp := true
	p.isFirstBlob = &tmp
	p.data.oldBlobStateRoot = ZeroHash
	p.data.oldBlobAccInputHash = ZeroHash
	p.data.oldNumBlob = 0
	p.data.oldStateRoot = ZeroHash
}

func (p *ProcessBlobInnerProcessRequestBuilder) setCurrentSequence(seq BlobSequence) {
	p.data.coinbase = seq.L2Coinbase
}

func (p *ProcessBlobInnerProcessRequestBuilder) setPreviousSequence(previousSequence BlobSequence) {
	if p.isFirstBlob != nil && !*p.isFirstBlob {
		p.err = errors.New("the blob is the first blob, you can't set a previous blob")
	}
	tmp := false
	p.isFirstBlob = &tmp
	p.data.oldBlobAccInputHash = previousSequence.FinalAccInputHash
	p.data.oldNumBlob = previousSequence.LastBlobSequenced
}

func (p *ProcessBlobInnerProcessRequestBuilder) setBlob(blob *BlobInner) {
	if p.err != nil {
		return
	}
	p.data.lastL1InfoTreeIndex = blob.L1InfoLeafIndex
	p.data.lastL1InfoTreeRoot = blob.L1InfoTreeRoot
	p.data.timestampLimit = uint64(blob.MaxSequenceTimestamp.Unix()) // Convert time.Time to uint64
	p.data.zkGasLimit = blob.ZkGasLimit
}

func (p *ProcessBlobInnerProcessRequestBuilder) Build() (ProcessBlobInnerProcessRequest, error) {
	if p.err != nil {
		return ProcessBlobInnerProcessRequest{}, p.err
	}
	return p.data, nil
}

type ProcessBlobInnerResponse struct {
}

// ProcessBlobInner processes a blobInner and returns the splitted batches
func (s *State) ProcessBlobInner(ctx context.Context, request ProcessBlobInnerProcessRequest, data []byte) (*ProcessBlobInnerResponse, error) {
	requestExecutor := convertBlobInnerProcessRequestToExecutor(request, data)
	processResponse, err := s.executorClient.ProcessBlobInnerV3(ctx, requestExecutor)
	if err != nil {
		log.Errorf("Error processing blobInner: %v", err)
		return nil, err
	}
	return convertBlobInnerResponseToState(processResponse), nil
}

func convertBlobInnerResponseToState(response *executor.ProcessBlobInnerResponseV3) *ProcessBlobInnerResponse {
	return &ProcessBlobInnerResponse{}
}

func convertBlobInnerProcessRequestToExecutor(request ProcessBlobInnerProcessRequest, data []byte) *executor.ProcessBlobInnerRequestV3 {
	return &executor.ProcessBlobInnerRequestV3{
		OldBlobStateRoot:    request.oldBlobStateRoot.Bytes(),
		OldBlobAccInputHash: request.oldBlobAccInputHash.Bytes(),
		OldNumBlob:          request.oldNumBlob,
		OldStateRoot:        request.oldStateRoot.Bytes(),
		ForkId:              request.forkId,
		LastL1InfoTreeIndex: request.lastL1InfoTreeIndex,
		LastL1InfoTreeRoot:  request.lastL1InfoTreeRoot.Bytes(),
		TimestampLimit:      request.timestampLimit,
		Coinbase:            request.coinbase.String(),
		ZkGasLimit:          request.zkGasLimit,
		BlobType:            uint32(request.blobType),
		BlobData:            data,
	}
}
