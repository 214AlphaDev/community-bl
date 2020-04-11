package community_bl

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	vo "github.com/214alphadev/community-bl/value_objects"
	"time"
)

type communityService struct {
	memberRepository      MemberRepository
	applicationRepository ApplicationRepository
	onApplicationApproved []func(member MemberEntity)
}

func (s communityService) GetLastApplication(memberID MemberIdentifier, requesterID MemberIdentifier) (ApplicationEntity, error) {

	requester, err := s.memberRepository.FetchByID(requesterID)
	if err != nil {
		return ApplicationEntity{}, err
	}

	if requester == nil {
		return ApplicationEntity{}, errors.New("RequesterNotFound")
	}

	application, err := s.applicationRepository.FetchLast(memberID)
	if err != nil {
		return ApplicationEntity{}, err
	}

	if application == nil {
		return ApplicationEntity{}, errors.New("ApplicationNotFound")
	}

	if application.MemberID == requesterID {
		return *application, nil
	}

	if requester.Admin {
		return *application, nil
	}

	return ApplicationEntity{}, errors.New("NotAllowedToAccessApplication")

}

func (s communityService) Application(applicationID ApplicationID, requester MemberIdentifier) (ApplicationEntity, error) {

	member, err := s.memberRepository.FetchByID(requester)
	if err != nil {
		return ApplicationEntity{}, err
	}

	if member == nil {
		return ApplicationEntity{}, errors.New("MemberDoesNotExist")
	}

	if !member.Admin {
		return ApplicationEntity{}, errors.New("InsufficientPermissions")
	}

	application, err := s.applicationRepository.FetchByID(applicationID)
	if err != nil {
		return ApplicationEntity{}, err
	}

	if application == nil {
		return ApplicationEntity{}, errors.New("ApplicationDoesNotExist")
	}

	return *application, nil

}

func (s *communityService) Applications(query ApplicationsQuery, requester MemberIdentifier) ([]ApplicationEntity, error) {

	member, err := s.memberRepository.FetchByID(requester)
	if err != nil {
		return nil, err
	}

	if member == nil {
		return nil, errors.New("MemberDoesNotExist")
	}

	if !member.Admin {
		return nil, errors.New("InsufficientPermissions")
	}

	return s.applicationRepository.FetchByQuery(query)

}

func (s *communityService) ApplyForVerification(memberID MemberIdentifier, applicationText string) (ApplicationEntity, error) {

	fetchedApplication, err := s.applicationRepository.FetchLast(memberID)
	if err != nil {
		return ApplicationEntity{}, err
	}

	var apply = func() (ApplicationEntity, error) {

		application := ApplicationEntity{
			ID:              uuid.NewV4(),
			MemberID:        memberID,
			ApplicationText: applicationText,
			State:           ApplicationStatePending,
			CreatedAt:       time.Now(),
		}

		if err := s.applicationRepository.Save(application); err != nil {
			return ApplicationEntity{}, nil
		}

		return application, nil

	}

	switch fetchedApplication {
	case nil:
		return apply()
	default:
		switch fetchedApplication.State {
		case ApplicationStatePending:
			return ApplicationEntity{}, errors.New("PendingApplication")
		case ApplicationStateApproved:
			return ApplicationEntity{}, errors.New("AlreadyVerified")
		case ApplicationStateRejected:
			return apply()
		default:
			return ApplicationEntity{}, fmt.Errorf("application state: '%s' is invalid", fetchedApplication.State)
		}
	}

}

func (s *communityService) ApproveApplication(applicationID ApplicationID, reviewerID MemberIdentifier) error {

	reviewer, err := s.memberRepository.FetchByID(reviewerID)
	if err != nil {
		return err
	}

	if reviewer == nil {
		return fmt.Errorf("member that is supposed to reivew the application doesn't exist (member id: %s)", reviewerID.String())
	}

	if !reviewer.Admin {
		return errors.New("InsufficientPermissions")
	}

	application, err := s.applicationRepository.FetchByID(applicationID)
	if err != nil {
		return nil
	}

	if application == nil {
		return errors.New("ApplicationDoesNotExist")
	}

	member, err := s.memberRepository.FetchByID(application.MemberID)
	if err != nil {
		return err
	}

	if application.State != ApplicationStatePending {
		return errors.New("AlreadyReviewed")
	}

	application.ApprovedBy = &reviewer.ID
	now := time.Now()
	application.ApprovedAt = &now
	application.State = ApplicationStateApproved

	if err := s.applicationRepository.Save(*application); err != nil {
		return err
	}

	member.Verified = true

	if err := s.memberRepository.Save(*member); err != nil {
		return err
	}

	for _, onApproved := range s.onApplicationApproved {
		onApproved(*member)
	}

	return nil

}

func (s *communityService) RejectApplication(applicationID ApplicationID, reason string, reviewerID MemberIdentifier) error {

	reviewer, err := s.memberRepository.FetchByID(reviewerID)
	if err != nil {
		return err
	}

	if reviewer == nil {
		return errors.New("ReviewerDoesNotExist")
	}

	if !reviewer.Admin {
		return errors.New("InsufficientPermissions")
	}

	application, err := s.applicationRepository.FetchByID(applicationID)
	if err != nil {
		return err
	}

	if application == nil {
		return errors.New("ApplicationDoesNotExist")
	}

	if application.State != ApplicationStatePending {
		return errors.New("ApplicationReviewed")
	}

	application.RejectionReason = reason
	now := time.Now()
	application.RejectedAt = &now
	application.State = ApplicationStateRejected
	application.RejectedBy = &reviewer.ID

	return s.applicationRepository.Save(*application)

}

func (s communityService) Promote(emailAddress vo.EmailAddress) error {

	member, err := s.memberRepository.FetchByEmailAddress(emailAddress)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("couldn't find member by email address")
	}

	alreadyVerified := member.Verified

	member.Verified = true
	member.Admin = true

	if err := s.memberRepository.Save(*member); err != nil {
		return err
	}

	if !alreadyVerified {
		for _, onApproved := range s.onApplicationApproved {
			onApproved(*member)
		}
	}

	return nil

}

func (s *communityService) GetApplicationByID(id ApplicationID) (ApplicationEntity, error) {

	application, err := s.applicationRepository.FetchByID(id)

	switch err {
	case nil:
		return *application, nil
	default:
		return ApplicationEntity{}, err
	}

}

func (s *communityService) OnApplicationApproved(cb func(member MemberEntity)) {
	s.onApplicationApproved = append(s.onApplicationApproved, cb)
}
