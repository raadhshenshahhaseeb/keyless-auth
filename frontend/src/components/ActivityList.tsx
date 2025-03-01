import React from 'react';

interface ActivityListProps {
  selectedNetwork: {
    id: string;
    name: string;
  };
}

const ActivityList: React.FC<ActivityListProps> = ({ selectedNetwork }) => {
  return (
    <div className="activity-list">
      {/* Activity items will use selectedNetwork.id for filtering transactions */}
      <div className="activity-item">
        <div className="activity-network">{selectedNetwork.name}</div>
        {/* Add activity details here */}
      </div>
    </div>
  );
};

export default ActivityList;